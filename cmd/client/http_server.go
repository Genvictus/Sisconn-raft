package main

import (
	"encoding/json"

	"log"
	"net/http"
)

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

func LoggingMiddleware(next httpHandlerFunc) httpHandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next(w, r)

        log.Println(
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            r.Form,
        )
    }
}

func handlePing(w http.ResponseWriter, r *http.Request) {
    // TODO
    message, _ := json.Marshal(Ping())
    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    // TODO
}

func handleSet(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")
    value := r.Form.Get("value")

    message, _ := json.Marshal(Set([]string{key, value}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleStrln(w http.ResponseWriter, r *http.Request) {
    // TODO
}

func handleDel(w http.ResponseWriter, r *http.Request) {
    // TODO
}

func handleAppend(w http.ResponseWriter, r *http.Request) {
    // TODO
}

func handleRequestLog(w http.ResponseWriter, r *http.Request) {
    // TODO
}
