package winton

type Task struct {
	Task_UID   string `json:"uid"`
	Beacon_UID string `json:"beacon_uid"`
	Cmd        string `json:"cmd"`
	Status     string `json:"status"`
	Result     string `json:"result"`
}

func NewTask(task_uid string, beacon_uid string, cmd string) *Task {
	return &Task{
		Task_UID:   task_uid,
		Beacon_UID: beacon_uid,
		Cmd:        cmd,
		Status:     "pending",
		Result:     "",
	}
}

func (t *Task) UpdateTask(status string, result string) *Task {
	t.Status = status
	t.Result = result
	return t
}

func (t *Task) GetResult(client *Client) (string, int, error) {
	b64_result, err := client.Get_Response(t.Task_UID)
	if err != nil {
		return "", 0, err
	}
	if len(b64_result.Results) > 0 {
		for _, result := range b64_result.Results {
			res := result.Result
			size := len([]byte(res))
			return res, size, nil
		}
	}

	return "", 0, nil
}
