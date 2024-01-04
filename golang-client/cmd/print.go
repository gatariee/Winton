package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

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
        tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
    )
    table.SetRowSeparator("-")
    table.SetCenterSeparator("+")
    table.SetColumnSeparator("|")

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

func (p *Print) BeaconSent(x int) {
	p.Infof("Tasked beacon, sent %d bytes.", x)
}
