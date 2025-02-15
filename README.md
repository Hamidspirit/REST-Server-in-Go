# About
Implementing REST servers in GO via different aproaches
[refrence](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/ "This entire website is cool")

## REST with std lib

Server is the classic Todo app. JSON data encoding.
Datalayer will contain a simple in memory database and logic.

task: text , tags , due

- **text**: string

- **tags**: []string

- **due**: time object

> This is also good resource for servers in Go and multiplexer [mux](https://dev.to/jpoly1219/what-even-is-a-mux-4fng)

## FrameWork (Gin)

Frame work will help to write less code and you will enjoy the convenience.

datalayer is the same, but i can see difference in handlers.

i can find collection of gin middleware on [here] (https://github.com/gin-contrib/)

## 
