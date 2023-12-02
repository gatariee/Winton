package handler

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