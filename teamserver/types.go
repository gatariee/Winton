package main

type Agent struct {
	IP       string
	ExtIP    string
	Hostname string
	Sleep    string
	Jitter   string
	OS       string
	UID      string
	PID      string
}

type Task struct {
	UID       string
	CommandID string
	Command   string
}

type CommandData struct {
	CommandID string
	Command   string
}

type Result struct {
	CommandID string
	Result    string
}

type Callback struct {
	AgentUID     string
	LastCallback int
}