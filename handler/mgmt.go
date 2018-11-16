package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"fmt"
	"noname/models"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func MgmtHandler(w http.ResponseWriter, r *http.Request) {

	method := r.Method
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])
	var content []byte

	switch method {
	case GET:
		content = mgmntGet()
	case POST:
		mgmntCreate(r)
	case PUT:
		mgmntUpdate()
	case DELETE:
		mgmntDelete(flagName)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//js, err := json.Marshal(method)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(content)
}

func mgmntGet() []byte {

	flags := models.RedisClient.Keys("flag-*")
	// scan the []interface{} slice into a []int slice
	keys, err := flags.Result()
	if err != nil {
		log.Fatal(err)
	}

	var flagArr []models.Flag

	for _, val := range keys {
		flagValues := models.RedisClient.HGetAll(val).Val()
		fmt.Println("name: " + val)
		fmt.Println("status: " + flagValues["status"])
		fmt.Println("ratio: " + flagValues["ratio"])
		fmt.Println("sticky: " + flagValues["sticky"])
		ratio, err := strconv.Atoi(flagValues["ratio"])
		if err != nil {
		}
		flag := models.Flag{val, flagValues["status"], flagValues["sticky"], ratio}
		flagArr = append(flagArr, flag)
	}

	js, err := json.Marshal(flagArr)
	if err == nil {
	}
	return js

}

func mgmntCreate(r *http.Request) {
	fmt.Println(r.Body)

	name := r.FormValue("name")
	status := r.FormValue("status")
	sticky := r.FormValue("sticky")
	ratio := r.FormValue("ratio")

	ratioValue, err := strconv.Atoi(ratio)
	if err != nil {
	}

	flag := models.Flag{name, status, sticky, ratioValue}
	writeFlag(flag)
}

func mgmntUpdate() {

}

func mgmntDelete(flagName string) {
	flagVar := "flag-" + flagName
	models.RedisClient.Del(flagVar)
}

func writeFlag(flag models.Flag) {

	// var fields map[string]interface{}
	var fields map[string]interface{}
	fields = make(map[string]interface{})

	fields["ratio"] = strconv.Itoa(flag.Ratio)
	fields["status"] = flag.Status
	fields["sticky"] = flag.Sticky
	flag.Name = "flag-" + strings.ToLower(flag.Name)

	models.RedisClient.HMSet(flag.Name, fields)
	//fmt.Println(result)
}
