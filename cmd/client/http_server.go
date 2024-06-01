package main

import (
	"encoding/json"

	"log"
	"net/http"
)

type httpHandlerFunc func(w http.ResponseWriter, r *http.Request)

func loggingMiddleware(next httpHandlerFunc) httpHandlerFunc {
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
    message, _ := json.Marshal(ping())
    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")

    message, _ := json.Marshal(get([]string{key}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleSet(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")
    value := r.Form.Get("value")

    message, _ := json.Marshal(set([]string{key, value}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleStrln(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")

    message, _ := json.Marshal(strln([]string{key}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleDel(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")

    message, _ := json.Marshal(del([]string{key}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleAppend(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data", http.StatusBadRequest)
        return
    }

    key := r.Form.Get("key")
    value := r.Form.Get("value")

    message, _ := json.Marshal(Append([]string{key, value}))

    w.WriteHeader(http.StatusOK)
    w.Write(message)
}

func handleRequestLog(w http.ResponseWriter, r *http.Request) {
    message, _ := json.Marshal(requestLog())
    w.WriteHeader(http.StatusOK)
    w.Write(message)
}
