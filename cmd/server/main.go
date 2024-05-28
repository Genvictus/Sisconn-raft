package main

import (
	"Sisconn-raft/raft/transport"
	"flag"
	"fmt"
)

// ServerParserInfo holds server command line arguments
type ServerParserInfo struct {
	Host string
	Port int
}

func main() {
	var serverInfo ServerParserInfo
	flag.StringVar(&serverInfo.Host, "host", "localhost", "Server hostname that will be used by client, default=localhost")
	flag.IntVar(&serverInfo.Port, "port", 6969, "Server port that will be used by client, default=6969")
	flag.Parse()

	server_address := transport.NewAddress(serverInfo.Host, serverInfo.Port)

	fmt.Println("Server started at", &server_address)
	
	// TODO: start server logic
}
