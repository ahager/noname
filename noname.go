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
	Status string
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])

	// debug
	flag, err := checkFlag(flagName)

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println(string(js))
}

func checkFlag(flagName string) (Flag, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	var flag Flag

	flagVar := "flag-" + flagName

	redisValue, err := client.Get(flagVar).Result()
	if err == redis.Nil {
		fmt.Println("flag does not exist")
	} else if err != nil {
		panic(err)
	} else {
		flag = Flag{flagVar, redisValue}
	}
	return flag, err
}

func main() {
	fmt.Println("Listen on Port 8080")

	http.HandleFunc("/flag/", handler)
	http.ListenAndServe(":8080", nil)
}
