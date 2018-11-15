package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "strings"
  "log"
)

// create a flag
// set a flag to active
// set a flag to inactive


type Flag struct {
  Name   string
  Status bool
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])

	// debug
	flag := Flag{flagName, true}

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println(string(js))
}

func main() {
	fmt.Println("Listen on Port 8080")

    http.HandleFunc("/flag/", handler)
    http.ListenAndServe(":8080", nil)
}
