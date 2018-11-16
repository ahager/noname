package handler

import (
    "encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

    "noname/models"
)

func FlagHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	flagName := string(path[2])

	clientID := string(path[3])
	if clientID == "" {
		clientID = models.CreateClientId(r)
	}

	var forcedState string
	if len(path) == 5 {
		forcedState = string(path[4])
	}

    num := rand.Intn(100)

	flag, err := checkFlag(flagName, clientID)
    if err != nil {
		log.Fatal(err)
	}

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
}

func checkFlag(flagName string, clientId string) (models.Flag, error) {

	var flag models.Flag

	flagVar := "flag-" + flagName
	redisMap, err := models.RedisClient.HGetAll(flagVar).Result()

	if err != nil {
		panic(err)
	} else {
		ratio, err := strconv.Atoi(redisMap["ratio"])
		if err != nil {
            panic(err)
        }
		flag = models.Flag{flagName, redisMap["status"], clientId, redisMap["sticky"], ratio}
	}

	return flag, err
}
