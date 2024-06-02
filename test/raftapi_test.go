package test

import (
	"Sisconn-raft/raft"
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"context"
	"log"
	"net"
	"os"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client pb.RaftServiceClient
)

func TestMain(m *testing.M) {
	serverAddress := transport.NewAddress("localhost", 1234)
	raftNode := raft.NewNode(serverAddress.String())
	lis, err := net.Listen("tcp", serverAddress.String())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	go startGRPCServer(raftNode, lis)

	conn, err := grpc.NewClient(serverAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	client = pb.NewRaftServiceClient(conn)

	code := m.Run()

	os.Exit(code)
}

func TestPing(t *testing.T) {
	ctx := context.Background()
	request := &pb.PingRequest{}
	response, err := client.Ping(ctx, request)
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
	expected := "OK"
	if response.Response != expected {
		t.Errorf("Expected response: %s, but got: %s", expected, response.Response)
	}
}

func TestReqLog(t *testing.T) {
	// Set up
	ctx := context.Background()
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	setRequest1 := &pb.KeyValuedRequest{Key: key1, Value: value1}
	setRequest2 := &pb.KeyValuedRequest{Key: key2, Value: value2}
	_, err := client.Set(ctx, setRequest1)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	_, err = client.Set(ctx, setRequest2)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Test ReqLog
	reqLogRequest := &pb.LogRequest{}
	reqLogResponse, err := client.ReqLog(ctx, reqLogRequest)
	if err != nil {
		t.Fatalf("ReqLog failed: %v", err)
	}

	// Verify ReqLog
	expectedLogEntries := []*pb.LogEntry{
		{Key: key1, Value: value1},
		{Key: key2, Value: value2},
	}
	if !reflect.DeepEqual(reqLogResponse.LogEntries, expectedLogEntries) {
		t.Errorf("Expected log entries: %v, but got: %v", expectedLogEntries, reqLogResponse.LogEntries)
	}
}

func TestCommit(t *testing.T) {
	ctx := context.Background()

	// TODO proper commit test
	request := &pb.CommitRequest{}
	response, err := client.Commit(ctx, request)
	if err != nil {
		t.Errorf("Commit failed: %v", err)
	}
	expected := "OK"
	if response.Response != expected {
		t.Errorf("Expected response: %s, but got: %s", expected, response.Response)
	}
}

func TestSetAndGet(t *testing.T) {
	// Set
	ctx := context.Background()
	key := "exampleKey"
	value := "exampleValue"
	setRequest := &pb.KeyValuedRequest{Key: key, Value: value}
	setResponse, err := client.Set(ctx, setRequest)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	expectedSetResponse := "OK"
	if setResponse.Response != expectedSetResponse {
		t.Errorf("Expected Set response: %s, but got: %s", expectedSetResponse, setResponse.Response)
	}

	// Get
	getRequest := &pb.KeyedRequest{Key: key}
	getResponse, err := client.Get(ctx, getRequest)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	expectedGetValue := value
	if getResponse.Value != expectedGetValue {
		t.Errorf("Expected value: %s, but got: %s", expectedGetValue, getResponse.Value)
	}
}

func TestDel(t *testing.T) {
	// Set
	ctx := context.Background()
	key := "exampleKey"
	value := "exampleValue"
	setRequest := &pb.KeyValuedRequest{Key: key, Value: value}
	_, err := client.Set(ctx, setRequest)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Del
	delRequest := &pb.KeyedRequest{Key: key}
	delResponse, err := client.Del(ctx, delRequest)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	expectedValue := value
	if delResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, delResponse.Value)
	}

	// Get
	getRequest := &pb.KeyedRequest{Key: key}
	getResponse, err := client.Get(ctx, getRequest)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	expectedValue = ""
	if getResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, getResponse.Value)
	}
}

func TestAppend(t *testing.T) {
	// Set
	ctx := context.Background()
	key := "exampleKey"
	initialValue := "initialValue"
	setRequest := &pb.KeyValuedRequest{Key: key, Value: initialValue}
	_, err := client.Set(ctx, setRequest)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Append
	appendValue := " appended"
	appendRequest := &pb.KeyValuedRequest{Key: key, Value: appendValue}
	appendResponse, err := client.Append(ctx, appendRequest)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	expectedResponse := "OK"
	if appendResponse.Response != expectedResponse {
		t.Errorf("Expected response: %s, but got: %s", expectedResponse, appendResponse.Response)
	}

	getRequest := &pb.KeyedRequest{Key: key}
	getResponse, err := client.Get(ctx, getRequest)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	expectedValue := initialValue + appendValue
	if getResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, getResponse.Value)
	}
}

func startGRPCServer(node *raft.RaftNode, lis net.Listener) {
	grpcServer := grpc.NewServer()
	pb.RegisterRaftServiceServer(grpcServer, &raft.ServiceServer{Server: node})
	pb.RegisterRaftServer(grpcServer, &raft.RaftServer{Server: node})
	grpcServer.Serve(lis)
}
