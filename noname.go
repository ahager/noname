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
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type Flag struct {
	Name     string
	Status   string
	ClientId string
	Sticky   string
	Ratio    int
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

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

func mgmntHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	method := r.Method
	w.Header().Set("Content-Type", "application/json")

	switch method {
	case GET:
		mgmntGet()
	case POST:
		mgmntCreate()
	case PUT:
		mgmntUpdate()
	case DELETE:
		mgmntDelete()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	js, err := json.Marshal(method)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(js)

}

func handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])

	clientID := string(path[3])
	if clientID == "" {
		clientID = createClientID(r)
	}

	var forcedState string
	if len(path) == 5 {
		forcedState = string(path[4])
	}

	flag, err := checkFlag(flagName, clientID)

	num := rand.Intn(100)

	//fmt.Println(flag.Status)
	//fmt.Println(num)

	if flag.Status == "1" {
		if forcedState != "" {
			flag.Status = forcedState
		} else {
			if num > flag.Ratio {
				flag.Status = "0"
			}
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


func writeFlag(flag Flag) {

    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // var fields map[string]interface{}
    var fields map[string]interface{}
    fields = make(map[string]interface{})

    fields["Ratio"] = strconv.Itoa(flag.Ratio)
    fields["Status"] = flag.Status
    fields["Sticky"] = flag.Sticky

    result := client.HMSet(flag.Name, fields)
    fmt.Println(result)
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
		if err != nil {
		}
		flag = Flag{flagName, redisMap["status"], clientId, redisMap["sticky"], ratio}
	}

	return flag, err
}

// var client redis.Client

func main() {

    /* client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    }) */

    flag := Flag{ "Oida", "1", "", "1", 75 }
    writeFlag(flag)

	fmt.Println("Listen on Port 8080")

	http.HandleFunc("/flag/", handler)
	http.HandleFunc("/mgmnt/", mgmntHandler)

	http.ListenAndServe(":8080", nil)
}

func mgmntGet() {

}

func mgmntCreate() {

}

func mgmntUpdate() {

}

func mgmntDelete() {

}
