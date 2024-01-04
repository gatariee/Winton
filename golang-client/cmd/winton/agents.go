package winton

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NewClient(Teamserver string) *Client {
	return &Client{
		Teamserver: Teamserver,
	}
}

func (c *Client) GetAgentList() ([]Agent, error) {
	if c.Teamserver[:4] != "http" {
		c.Teamserver = "http://" + c.Teamserver
	}

	resp, err := http.Get(c.Teamserver + "/agents")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)

	var data struct {
		Agents []Agent `json:"agents"`
	}

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data.Agents, nil
}

func (c *Client) FindAgentByUID(uid string) (Agent, bool) {
	for _, agent := range c.AgentList {
		if agent.UID == uid {
			return agent, true
		}
	}
	return Agent{}, false
}
