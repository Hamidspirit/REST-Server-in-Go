package main

import (
  "net/http"
  "fmt"
  "log"
)
func main() {
  fmt.Println("Starting a server on port 8080...:")
  err = http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Fatal("Error while running the server")
  }
}
 
 
