package models

import (
    "crypto/md5"
	"encoding/hex"
    "io"
    "net/http"
)

type Client struct {
	Id         string
	Agent      string
	LastAccess string
}

func WriteClientLog(r *http.Request) {

}

func CreateClientId(r *http.Request) string {
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
