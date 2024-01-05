package winton

type Client struct {
	AgentList     []Agent
	Tasks         []Task
	AgentID       string
	AgentHostname string
	AgentIP       string
	BeaconSleep   string
	Teamserver    string
}

type Agent struct {
	IP       string `json:"IP"`
	ExtIP    string `json:"ExtIP"`
	Hostname string `json:"Hostname"`
	Sleep    string `json:"Sleep"`
	Jitter   string `json:"Jitter"`
	OS       string `json:"OS"`
	UID      string `json:"UID"`
	PID      string `json:"PID"`
}

type File struct {
	Filename string `json:"Filename"`
	Size     int    `json:"Size"`
	IsDir    bool   `json:"IsDir"`
	ModTime  string `json:"ModTime"`
}

type CommandData struct {
	CommandID string `json:"CommandID"`
	Command   string `json:"Command"`
}

type Result struct {
	CommandID string `json:"CommandID"`
	Result    string `json:"Result"`
}

type ResultList struct {
	Results []Result `json:"Results"`
}

type ResponseData struct {
	Message string `json:"message"`
	Task_ID string `json:"uid"`
}
