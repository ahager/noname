package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
    "strconv"

	"github.com/go-redis/redis"
)

type Flag struct {
	Name     string
	Status   string
	ClientId string
	Sticky   string
	Ratio    int
}

func createClientID(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	acceptLanguage := r.Header.Get("Accept-Language")
	// remoteAddr := r.Header.Get("X-FORWARDED-FOR")

	h := md5.New()
	io.WriteString(h, remoteAddr)
	io.WriteString(h, userAgent)
	io.WriteString(h, acceptLanguage)

	return hex.EncodeToString(h.Sum(nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])

	clientID := string(path[3])
	if clientID == "" {
		clientID = createClientID(r)
	}

	flag, err := checkFlag(flagName, clientID)

	num := rand.Intn(100)

	fmt.Println(flag.Status)
	fmt.Println(num)

	if flag.Status == "1" {
		if num > flag.Ratio {
			flag.Status = "0"
		}
	}

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	//fmt.Println(clientID)
}

func checkFlag(flagName string, clientId string) (Flag, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//pong, err := client.Ping().Result()
	//fmt.Println(pong, err)
	var flag Flag

	flagVar := "flag-" + flagName
    redisMap, err := client.HGetAll(flagVar).Result()

	if err == redis.Nil {
		fmt.Println("flag does not exist")
	} else if err != nil {
		panic(err)
	} else {
        ratio, err := strconv.Atoi(redisMap["ratio"])
        if err != nil {}
        flag = Flag{flagName, redisMap["status"], clientId, redisMap["sticky"], ratio}
	}

	return flag, err
}

func main() {
	fmt.Println("Listen on Port 8080")

	http.HandleFunc("/flag/", handler)
	http.ListenAndServe(":8080", nil)
}
