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

    	switch method {
    	case GET:
    		mgmntGet()
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

    	js, err := json.Marshal(method)
    	if err != nil {
    		log.Fatal(err)
    	}
    	w.Write(js)
}

func mgmntGet() {

}

func mgmntCreate(r *http.Request) {
	fmt.Println(r.Body)

	name := r.FormValue("Name")
	status := r.FormValue("Status")
	sticky := r.FormValue("Sticky")
	ratio := r.FormValue("Ratio")

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

	fields["Ratio"] = strconv.Itoa(flag.Ratio)
	fields["Status"] = flag.Status
	fields["Sticky"] = flag.Sticky
	flag.Name = "flag-" + strings.ToLower(flag.Name)

	models.RedisClient.HMSet(flag.Name, fields)
	//fmt.Println(result)
}
