package main

import (
    "fmt"
    "net/http"
    "noname/handler"
    "noname/models"
)

func main() {

    models.InitRedis("localhost:6379", "", 0)

	http.HandleFunc("/flag/", handler.FlagHandler)
	http.HandleFunc("/mgmnt/", handler.MgmtHandler)

    fmt.Println("Listen on Port 8080")
	http.ListenAndServe(":8080", nil)
}
