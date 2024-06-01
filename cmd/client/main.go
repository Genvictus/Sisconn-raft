package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type ClientParserInfo struct {
	ServerHost string
	ServerPort int
}

var (
	targetServer        *t.Address
	serviceClient       pb.RaftServiceClient
	conn                *grpc.ClientConn
	err                 error
	ClientLogger        *log.Logger
	ClientVerboseLogger *log.Logger
)

func main() {
	var clientInfo ClientParserInfo
	flag.StringVar(&clientInfo.ServerHost, "host", "localhost", "Server host for target connection, default=localhost")
	flag.IntVar(&clientInfo.ServerPort, "port", 6969, "Server port for target connection, default=6969")
	flag.Parse()

	serverAddress := t.NewAddress(clientInfo.ServerHost, clientInfo.ServerPort)

	ClientLogger = log.New(os.Stdout, "[Raft] Client: ", 0)

	t.LogPrint(ClientLogger, "Client started")
	t.LogPrint(ClientLogger, fmt.Sprintf("Connecting to server at %s", serverAddress.String()))
	// fmt.Println("Client Started")
	// fmt.Println("Connecting to server at", &serverAddress)

	SetTargetServer(&serverAddress)
	defer conn.Close()

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
