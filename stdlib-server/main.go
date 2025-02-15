//POST   /task/              :  create a task, returns ID
// GET    /task/<taskid>      :  returns a single task by ID
// GET    /task/              :  returns all tasks
// DELETE /task/<taskid>      :  delete a task by ID
// GET    /tag/<tagname>      :  returns list of tasks with this tag
// GET    /due/<yy>/<mm>/<dd> :  returns list of tasks due by this date

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"something.some/datalayer"
	"something.some/middleware"
)

type taskServer struct {
	store *datalayer.TaskStore
}

func NewTaskServer() *taskServer {
	store := datalayer.New()
	return &taskServer{store: store}
}

func (ts *taskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handl task create at %s \n", r.URL.Path)

	// This struct will help to deserialize the request
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	// This struct will help to return the id from CreateTask function
	type ResponseId struct {
		ID int `json:"id"`
	}

	// Make sure media type is JSON
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expected application/json content-type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var reqtask RequestTask // Initiate RequestTask
	// Decode JSON data from Reader to reqtask
	if err := dec.Decode(&reqtask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create task and return task id
	id := ts.store.CreateTask(reqtask.Text, reqtask.Tags, reqtask.Due)
	js, err := json.Marshal(ResponseId{ID: id}) // Marshal id into JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js) // Return id
}

func (ts *taskServer) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle get all task at %s", r.URL.Path)

	allTask := ts.store.GetAllTask()
	js, err := json.Marshal(allTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle getting task by id at : %s", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle getting task by tag at: %s", r.URL.Path)

	tag := r.PathValue("tag")

	tasks := ts.store.GetTaskByTags(tag)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getByDueDate(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle getting task by due date at: %s", r.URL.Path)

	badRequest := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", r.URL.Path), http.StatusBadRequest)
	}
	year, errYear := strconv.Atoi(r.PathValue("year"))
	month, errMonth := strconv.Atoi(r.PathValue("month"))
	day, errDay := strconv.Atoi(r.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil {
		badRequest()
		return
	}

	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)

	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle deleting all tasks at: %s", r.URL.Path)

	ts.store.DeleteAll()

}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle deleting a task based on id at: %s", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func main() {
	fmt.Println("hi mom")
	mux := http.NewServeMux()
	server := NewTaskServer()

	mux.HandleFunc("POST /task/", server.CreateTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTaskHandler)
	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	mux.HandleFunc("GET /tag/{tag}/", server.tagHandler)
	mux.HandleFunc("GET /task/{year}/{month}/{day}/", server.getByDueDate)
	mux.HandleFunc("DELETE /task/", server.deleteAllHandler)
	mux.HandleFunc("POST /task/{id}", server.deleteTaskHandler)

	address := "localhost:4567"

	handler := middleware.Logging(mux)
	handler = middleware.PanicRecovery(handler)

	fmt.Printf("Server started on : %s", address)
	log.Fatal(http.ListenAndServe(address, handler))
}
