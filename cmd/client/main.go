package main

import (
	"Sisconn-raft/raft/transport"
	"flag"
	"fmt"
)

type ClientParserInfo struct {
	ServerHost string
	ServerPort int
}

func main() {
	var clientInfo ClientParserInfo
	flag.StringVar(&clientInfo.ServerHost, "host", "localhost", "Server host for target connection, default=localhost")
	flag.IntVar(&clientInfo.ServerPort, "port", 6969, "Server port for target connection, default=6969")
	flag.Parse()

	server_address := transport.NewAddress(clientInfo.ServerHost, clientInfo.ServerPort)

	fmt.Println("Client Started")
	fmt.Println("Connecting to server at", &server_address)

	// TODO: client execute command logic
}
