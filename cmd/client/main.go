package main

import (
	"Sisconn-raft/raft/transport"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

	RunCommandLoop()
}

func RunCommandLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		ExecuteCommand(input)
	}
}
