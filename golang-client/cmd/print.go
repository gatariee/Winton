package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	winton "cli/cmd/winton"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type Print struct {
	errorPrefix string
	okPrefix    string
	infoPrefix  string
}

func NewPrint() *Print {
	return &Print{
		errorPrefix: "[!] ",
		okPrefix:    "[+] ",
		infoPrefix:  "[*] ",
	}
}

func (p *Print) Welcome(operator string) {
	operator_c := color.New(color.FgWhite, color.Bold)
	dt := color.New(color.FgMagenta, color.Bold)
	currentTime := time.Now().Format("02/01/2006 03:04:05 PM")

	fmt.Printf("%s - ", dt.Sprintf(currentTime))
	fmt.Printf("Welcome back, %s.", operator_c.Sprintf(operator))

	p.Linebreak()
}

func (p *Print) Linebreak() {
	fmt.Println()
}

func (p *Print) AgentsTable(agents []winton.Agent) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hostname", "IP Address", "Operating System", "Sleep Time", "Jitter", "Process ID", "UID"})

	for _, agent := range agents {
		row := []string{
			agent.Hostname,
			agent.IP,
			agent.OS,
			agent.Sleep,
			agent.Jitter,
			agent.PID,
			agent.UID,
		}
		table.Append(row)
	}

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
	)

	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

	table.Render()
}

func (p *Print) ConfigTable(config map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Setting", "Value"})

	for key, value := range config {
		table.Append([]string{key, value})
	}

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})

	table.Render()
}

func (p *Print) ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (p *Print) TasksTable(tasks []winton.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Task UID", "Beacon UID", "Command", "Status", "Result (Base64-encoded)"})

	for _, task := range tasks {
		row := []string{
			task.Task_UID,
			task.Beacon_UID,
			task.Cmd,
			task.Status,
			task.Result,
		}

		switch task.Status {
		case "complete":
			row[3] = color.New(color.FgHiGreen).Sprint(task.Status)
		case "failed":
			row[3] = color.New(color.FgHiRed).Sprint(task.Status)
		default:
			row[3] = color.New(color.FgHiYellow).Sprint(task.Status)
		}

		table.Append(row)
	}

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
	)

    table.SetRowSeparator(" ")
    table.SetCenterSeparator(" ")
    table.SetColumnSeparator(" ")
	table.Render()
}

func (p *Print) Errorf(format string, a ...interface{}) {
	fmt.Print(color.RedString(p.errorPrefix))
	fmt.Printf(format+"\n", a...)
}

func (p *Print) Okf(format string, a ...interface{}) {
	fmt.Print(color.GreenString(p.okPrefix))
	fmt.Printf(format+"\n", a...)
}

func (p *Print) Infof(format string, a ...interface{}) {
	fmt.Print(color.CyanString(p.infoPrefix))
	fmt.Printf(format+"\n", a...)
}

func (p *Print) BeaconRecv(x int) {
	p.Infof("Beacon called home, received %d bytes.", x)
}

func (p *Print) BeaconSent(x int, uid string, action string) {
	uidColor := color.New(color.FgHiBlue)
	actionColor := color.New(color.FgHiGreen)

	message := fmt.Sprintf("Tasked beacon [%s] to %s, sent %d bytes.", uidColor.Sprint(uid), actionColor.Sprint(action), x)
	p.Infof(message)
}
