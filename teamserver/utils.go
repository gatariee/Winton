package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func _winton_print(data string) {
	fmt.Println("[*] " + data)
}

func _winton_error(data string) {
	fmt.Println("[!] " + data)
}

func _winton_usage() {
	fmt.Println("Usage: ./teamserver <ip> <port> <password>")
}


func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (ts *TeamServer) checkBeacons() {
	for {
		time.Sleep(1 * time.Second)
		for i, callback := range ts.AgentCallbacks {
			agent, found := ts.findAgentByUID(callback.AgentUID)
			if !found {
				continue
			}
			ts.AgentCallbacks[i].LastCallback++
			agentSleep, _ := strconv.Atoi(agent.Sleep)
			if ts.AgentCallbacks[i].LastCallback > agentSleep+5 {
				fmt.Printf("[!] Agent [%s] has gone offline.\n", agent.UID)
				ts.removeAgent(agent.UID)
			}
		}
	}
}

func (ts *TeamServer) findAgentByUID(uid string) (Agent, bool) {
	for _, agent := range ts.AgentList {
		if agent.UID == uid {
			return agent, true
		}
	}
	return Agent{}, false
}

func (ts *TeamServer) removeAgent(uid string) {
	for i, agent := range ts.AgentList {
		if agent.UID == uid {
			ts.AgentList = append(ts.AgentList[:i], ts.AgentList[i+1:]...)
			return
		}
	}
}