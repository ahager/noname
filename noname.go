package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
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
	flag := Flag{flagName, false}

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println(string(js))
}

func checkFlag(flagName string) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func main() {
	fmt.Println("Listen on Port 8080")

	http.HandleFunc("/f/", handler)
	http.ListenAndServe(":8080", nil)
}
