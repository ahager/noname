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
)

// create a flag
// set a flag to active
// set a flag to inactive


type Flag struct {
  Name     string
  Status   bool
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
        // fmt.Println(clientId)
    }

	// debug
	flag := Flag{flagName, false, clientId}

	js, err := json.Marshal(flag)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println(clientId)
}

func main() {
	fmt.Println("Listen on Port 8080")

    http.HandleFunc("/flag/", handler)
    http.ListenAndServe(":8080", nil)
}
