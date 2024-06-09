package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	message, _ := json.Marshal(requestLog(serviceClient))
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func handleGetAllNode(w http.ResponseWriter, r *http.Request) {
	message, _ := json.Marshal(GetAllNode([]string{}))

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func handleAddNode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	host := r.Form.Get("host")
	port := r.Form.Get("port")

	message, _ := json.Marshal(AddNode([]string{host, port}))

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func handleRemoveNode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	host := r.Form.Get("host")
	port := r.Form.Get("port")

	message, _ := json.Marshal(RemoveNode([]string{host, port}))

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func handleGetLog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	host := r.Form.Get("host")
	port := r.Form.Get("port")

	fmt.Println("host dan port", host, port)
	portInt, err := strconv.Atoi(port)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"error": "Invalid port"}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
	}

	newServer := t.NewAddress(host, portInt)
	newClientLogger := log.New(os.Stdout, "[Raft] Client "+newServer.String()+" : ", 0)

	var newConn *grpc.ClientConn

	newConn, err = grpc.NewClient(targetServer.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(t.LoggingInterceptorClient(newClientLogger)))
	if err != nil {
		return
	}

	newServiceClient := pb.NewRaftServiceClient(newConn)

	message, _ := json.Marshal(requestLog(newServiceClient))

	w.WriteHeader(http.StatusOK)
	w.Write(message)

	newConn.Close()
}
