package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "strings"
  "log"
  "crypto/md5"
  "encoding/hex"
  "io"
  "github.com/go-redis/redis"
)


type Flag struct {
  Name     string
  Status   string
  ClientId string
}

func createClientId(r *http.Request) string {
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

    clientId := string(path[3])
    if clientId == "" {
        clientId = createClientId(r)
    }

	flag, err := checkFlag(flagName, clientId)

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println(clientId)
}

func checkFlag(flagName string, clientId string) (Flag, error) {
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
		flag = Flag{flagVar, redisValue, clientId}
	}
	return flag, err
}

func main() {
	fmt.Println("Listen on Port 8080")

	http.HandleFunc("/flag/", handler)
	http.ListenAndServe(":8080", nil)
}
