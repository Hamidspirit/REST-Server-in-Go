package datalayer

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	sync.Mutex
	// makes sure that only one goroutine can access at a time
	// this will give me the ability to lock or unlock tasks map and nextId field
	// this will prevent race conditions

	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	taskStore := &TaskStore{}
	taskStore.tasks = make(map[int]Task)
	taskStore.nextId = 0

	return taskStore
}

// Creates task and returns the id for that task
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()         // Lock mutex to ensure safe access to shared resource
	defer ts.Unlock() // Make sure to release the lock when function returns

	task := Task{
		ID:   ts.nextId, // Assign the current value of nextId to task ID
		Text: text,
		Due:  due,
	}

	task.Tags = make([]string, len(tags)) // Create new slice to hold tags
	copy(task.Tags, tags)                 // Copy the tags from input slice into tasks tags slice

	ts.tasks[ts.nextId] = task // Store task in map using nextId as key
	ts.nextId++                // Increment nextId for next task
	return task.ID             // Return the ID of created task
}

// GetTask will return task by given id
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, Ok := ts.tasks[id] // Check if task is in map
	if Ok {
		return task, nil
	} else {
		return Task{}, fmt.Errorf("Task with id: %d , was not found", id)
	}
}

// DeleteTask removes a task based on id and returns an error if failed
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	// Check if task exists
	if _, Ok := ts.tasks[id]; !Ok {
		return fmt.Errorf("Task with id: %d was not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// DeleteAll tasks
func (ts *TaskStore) DeleteAll() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

// GetAllTask returns all the tasks in memory
func (ts *TaskStore) GetAllTask() []Task {
	ts.Lock()
	defer ts.Unlock()

	tasks := make([]Task, 0, len(ts.tasks)) // Make an slice of Task with 0 length and 10 capacity
	// in Go slices will be returned as reference so i shoul not use pointers that makes dangling pointer

	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// GetTaskByTags returns all the tasks with same tag
func (ts *TaskStore) GetTaskByTags(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

Taskloop:
	for _, task := range ts.tasks {
		for _, tasktag := range task.Tags {
			if tasktag == tag {
				tasks = append(tasks, task)
				continue Taskloop // Could also use break to break out of inner loop and continue the outer loop
			}
		}
	}

	return tasks
}

// GetTaskByDueDate returns tasks with given due date
func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
