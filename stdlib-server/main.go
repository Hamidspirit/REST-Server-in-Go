// POST   /task/              :  create a task, returns ID
// GET    /task/<taskid>      :  returns a single task by ID
// GET    /task/              :  returns all tasks
// DELETE /task/<taskid>      :  delete a task by ID
// GET    /tag/<tagname>      :  returns list of tasks with this tag
// GET    /due/<yy>/<mm>/<dd> :  returns list of tasks due by this date
// This will be endpoints to implement using standard lib

package main

import (
	"fmt"

	"something.some/datalayer"
)

func main() {
	datalayer.New()
	fmt.Println("hi mom")
}
