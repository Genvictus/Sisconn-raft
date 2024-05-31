package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	serverAddress := t.NewAddress(clientInfo.ServerHost, clientInfo.ServerPort)

	fmt.Println("Client Started")
	fmt.Println("Connecting to server at", &serverAddress)

	conn, err := grpc.NewClient(targetServer.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial(targetServer.String()) <- deprecated
	if err != nil {
		return
	}
	defer conn.Close()

	serviceClient = pb.NewRaftServiceClient(conn)

	RunCommandLoop(&serverAddress)
}

func RunCommandLoop(serverAddress *t.Address) {
	reader := bufio.NewReader(os.Stdin)
	setTargetServer(serverAddress)

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
