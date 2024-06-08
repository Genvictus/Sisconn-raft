package main

import (
	"Sisconn-raft/raft"
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"flag"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// ServerParserInfo holds server command line arguments
type ServerParserInfo struct {
	Host string
	Port int
}

var ServerLogger *log.Logger
var ServerVerboseLogger *log.Logger
var RaftNode *raft.RaftNode

func main() {
	ServerLogger = log.New(os.Stdout, "[Raft] Server: ", 0)

	err := godotenv.Load()
	if err != nil {
		ServerLogger.Println("Error loading .env file. Using default values.")
	} else {
		raft.LoadRaftConfig()
	}

	var serverInfo ServerParserInfo
	flag.StringVar(&serverInfo.Host, "host", "localhost", "Server hostname that will be used by client, default=localhost")
	flag.IntVar(&serverInfo.Port, "port", 6969, "Server port that will be used by client, default=6969")
	flag.Parse()

	serverAddress := transport.NewAddress(serverInfo.Host, serverInfo.Port)

	transport.LogPrint(ServerLogger, "Server started")
	// fmt.Println("Server started at", &serverAddress)

	serverAddressStr := serverAddress.String()
	lis, err := net.Listen("tcp", serverAddressStr)

	if err != nil {
		ServerLogger.Fatalf("failed to listen on %s", serverAddressStr)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(transport.LoggingInterceptorServer(ServerLogger)),
	)

	RaftNode = raft.NewNode(serverAddressStr)
	RaftNode.AddConnections([]string{serverAddressStr})
	pb.RegisterRaftServiceServer(s, &raft.ServiceServer{Server: RaftNode})
	pb.RegisterRaftServer(s, &raft.RaftServer{Server: RaftNode})

	s.Serve(lis)
}
