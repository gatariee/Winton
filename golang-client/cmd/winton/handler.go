package winton

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Send_Task(task string, uid string) (ResponseData, error) {
	URL := c.Teamserver + "/tasks/" + uid
	commandData := CommandData{CommandID: "", Command: task}
	data, err := json.Marshal(commandData)
	if err != nil {
		return ResponseData{}, err
	}

	response, err := http.Post(URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return ResponseData{}, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var commandResponse ResponseData
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return ResponseData{}, err
		}

		err = json.Unmarshal(responseData, &commandResponse)
		if err != nil {
			return ResponseData{}, err
		}

		return commandResponse, nil
	}

	return ResponseData{}, fmt.Errorf("error sending task to agent")
}

func (c *Client) Get_Response(Task_ID string) (ResultList, error) {
	URL := c.Teamserver + "/results/" + Task_ID
	response, err := http.Get(URL)
	if err != nil {
		return ResultList{}, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var result ResultList
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return ResultList{}, err
		}

		err = json.Unmarshal(responseData, &result)
		if err != nil {
			return ResultList{}, err
		}

		return result, nil
	}

	return ResultList{}, fmt.Errorf("error getting response from agent")
}

func DecodeResult(result string) (string, error) {
	decodedResult, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return "", err
	}
	return string(decodedResult), nil
}
