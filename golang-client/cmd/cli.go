package cmd

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Session struct {
	currentSession string
	prompt         *readline.Instance
	print          *Print
}

func promptInit() string {
	username := viper.GetString("operator")
	dt := color.RedString(time.Now().Format("2006-01-02 15:04:05"))
	promptText := fmt.Sprintf("[%s] %s@winton ", dt, username)

	inputText := ">> "

	fullPrompt := promptText + inputText
	return fullPrompt
}

func NewSession() *Session {
	fullPrompt := promptInit()

	rl, err := readline.New(fullPrompt)
	if err != nil {
		panic(err)
	}

	p := NewPrint()

	return &Session{
		currentSession: "",
		prompt:         rl,
		print:          p,
	}
}

func (s *Session) Start() {
	defer s.prompt.Close()

	s.print.Welcome(viper.GetString("operator"))

	config := map[string]string{
		"Operator": viper.GetString("operator"),
		"Teamserver IP": viper.GetString("teamserver.ip"),
		"Teamserver Port": viper.GetString("teamserver.port"),
	}

	s.print.ConfigTable(config)

	for {

		update := promptInit()
		s.prompt.SetPrompt(update)
		s.print.Linebreak()

		line, err := s.prompt.Readline()
		if err != nil {
			break
		}

		if ok, err := s.processCommand(line); !ok {
			if err != nil {
				fmt.Println("error:", err)
			}
			break
		}
	}
}

func (s *Session) processCommand(command string) (bool, error) {
	switch command {
	case "exit":
		s.print.Infof("Exiting...")
		return false, nil

	case "test":
		s.print.BeaconRecv(123)
		return true, nil

	case "test2":
		s.print.BeaconSent(456)
		return true, nil

	case "clear":
		s.print.ClearScreen()
		return true, nil

	default:
		s.print.Infof("Unknown command: %s", command)
		return true, nil
	}
}
