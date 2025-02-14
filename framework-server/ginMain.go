package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/datalayer"
	"github.com/gin-gonic/gin"
)

type taskServer struct {
	store *datalayer.TaskStore
}

func NewTaskServer() *taskServer {
	store := datalayer.New()
	return &taskServer{store: store}
}

func (ts *taskServer) createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (ts *taskServer) getAllTaskHandler(c *gin.Context) {
	allTask := ts.store.GetAllTask()
	c.JSON(http.StatusOK, allTask)
}

func (ts *taskServer) deleteAllTaskHandler(c *gin.Context) {
	ts.store.DeleteAll()
}

func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ts *taskServer) deleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
}

func (ts *taskServer) getTaskByTagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")
	task := ts.store.GetTaskByTags(tag)
	c.JSON(http.StatusOK, task)
}

func (ts *taskServer) getTaskByDueDateHandler(c *gin.Context) {
	badRequestErr := func() {
		c.String(http.StatusBadRequest, "expected /due/<year>/<month>/<day>, got %s", c.FullPath())
	}

	year, err := strconv.Atoi(c.Params.ByName("year"))
	if err != nil {
		badRequestErr()
		return
	}

	month, err := strconv.Atoi(c.Params.ByName("month"))
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestErr()
		return
	}

	day, err := strconv.Atoi(c.Params.ByName("day"))
	if err != nil {
		badRequestErr()
		return
	}

	task := ts.store.GetTaskByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, task)
}

func main() {
	fmt.Println("hi mom")

	router := gin.Default()
	dataserv := NewTaskServer()

	router.POST("/task", dataserv.createTaskHandler)
	router.GET("/task", dataserv.getAllTaskHandler)
	router.DELETE("/task", dataserv.deleteAllTaskHandler)
	router.GET("/task/:id", dataserv.getTaskHandler)
	router.DELETE("/task/:id", dataserv.deleteTaskHandler)
	router.GET("/tag/:tag", dataserv.getTaskByTagHandler)
	router.GET("/due/:year/:month/:day", dataserv.getTaskByDueDateHandler)

	router.Run("localhost:4567")
}
