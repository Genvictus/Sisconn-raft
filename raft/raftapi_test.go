package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"context"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client pb.RaftServiceClient
)

func TestMain(m *testing.M) {
	serverAddress := transport.NewAddress("localhost", 1234)
	lis, err := net.Listen("tcp", serverAddress.String())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	raftNode := NewNode(serverAddress.String())
	raftNode.AddConnections([]string{serverAddress.String()})
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
		{},
		{Key: key1, Value: value1},
		{Key: key2, Value: value2},
	}
	if !reflect.DeepEqual(reqLogResponse.LogEntries, expectedLogEntries) {
		t.Errorf("Expected log entries: %v, but got: %v", expectedLogEntries, reqLogResponse.LogEntries)
	}
}

func TestCommit(t *testing.T) {
	ctx := context.Background()

	// Transaction list
	request := &pb.CommitRequest{
		CommitEntries: []*pb.CommitEntry{
			{Type: "set", Key: "CommitKey1", Value: "CommitValue1"},
			{Type: "set", Key: "CommitKey2", Value: "CommitValue2"},
			{Type: "set", Key: "CommitKey3", Value: "CommitValue3"},
			{Type: "set", Key: "CommitKey4", Value: "CommitValue4"},
			{Type: "del", Key: "CommitKey2", Value: ""},
			{Type: "append", Key: "CommitKey3", Value: " AppendValue3"},
			{Type: "set", Key: "CommitKey4", Value: "ResetValue4"},
		},
	}

	// Commit
	response, err := client.Commit(ctx, request)
	if err != nil {
		t.Errorf("Commit failed: %v", err)
	}
	expected := "OK (7 commands execution)"
	if response.Response != expected {
		t.Errorf("Expected response: %s, but got: %s", expected, response.Response)
	}

	// Check get value
	tests := []*pb.KeyValuedRequest{
		{Key: "CommitKey1", Value: "CommitValue1"},
		{Key: "CommitKey2", Value: ""},
		{Key: "CommitKey3", Value: "CommitValue3 AppendValue3"},
		{Key: "CommitKey4", Value: "ResetValue4"},
	}

	for _, tt := range tests {
		t.Run(tt.Key, func(t *testing.T) {
			getRequest := &pb.KeyedRequest{Key: tt.Key}
			getResponse, err := client.Get(ctx, getRequest)
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}

			if getResponse.Value != tt.Value {
				t.Errorf("Expected value for key %s: %s, but got: %s", tt.Key, tt.Value, getResponse.Value)
			}
		})
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
	key := "exampleKey2"
	value := "exampleValue2"
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

func TestStrln(t *testing.T) {
	// Set
	ctx := context.Background()
	key := "exampleKey3"
	value := "exampleValue3"
	setRequest := &pb.KeyValuedRequest{Key: key, Value: value}
	_, err := client.Set(ctx, setRequest)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Strln
	strlnRequest := &pb.KeyedRequest{Key: key}
	strlnResponse, err := client.Strln(ctx, strlnRequest)
	if err != nil {
		t.Fatalf("Strln failed: %v", err)
	}

	expectedValue := strconv.Itoa(len(value))
	if strlnResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, strlnResponse.Value)
	}

	// Append
	appendValue := " appended"
	appendRequest := &pb.KeyValuedRequest{Key: key, Value: appendValue}
	_, err = client.Append(ctx, appendRequest)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Strln
	strlnResponse, err = client.Strln(ctx, strlnRequest)
	if err != nil {
		t.Fatalf("Strln failed: %v", err)
	}

	expectedValue = strconv.Itoa(len(value) + len(appendValue))
	if strlnResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, strlnResponse.Value)
	}

	// Del
	delRequest := &pb.KeyedRequest{Key: key}
	_, err = client.Del(ctx, delRequest)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	strlnResponse, err = client.Strln(ctx, strlnRequest)
	if err != nil {
		t.Fatalf("Strln failed: %v", err)
	}

	expectedValue = "0"
	if strlnResponse.Value != expectedValue {
		t.Errorf("Expected value: %s, but got: %s", expectedValue, strlnResponse.Value)
	}
}

func TestRequestVote(t *testing.T) {

	// Server1 (Follower)
	serverAddress1 := transport.NewAddress("localhost", 2001)
	raftNode1 := NewNode(serverAddress1.String())

	raftServer1 := &RaftServer{Server: raftNode1}
	lis1, err := net.Listen("tcp", serverAddress1.String())
	if err != nil {
		t.Fatalf("Failed to listen for the first server: %v", err)
	}
	defer lis1.Close()
	grpcServer1 := grpc.NewServer()
	pb.RegisterRaftServer(grpcServer1, raftServer1)
	go grpcServer1.Serve(lis1)

	// Server2 (Candidate)
	serverAddress2 := transport.NewAddress("localhost", 2002)
	raftNode2 := NewNode(serverAddress2.String())
	raftServer2 := &RaftServer{Server: raftNode2}

	lis2, err := net.Listen("tcp", serverAddress2.String())
	if err != nil {
		t.Fatalf("Failed to listen for the second server: %v", err)
	}
	defer lis2.Close()
	grpcServer2 := grpc.NewServer()
	pb.RegisterRaftServer(grpcServer2, raftServer2)
	go grpcServer2.Serve(lis2)

	// Granted request
	request := &pb.RequestVoteArg{
		Term:         1,
		CandidateId:  "candidate1",
		LastLogTerm:  1,
		LastLogIndex: map[uint32]uint64{0: 0},
	}

	ctx := context.Background()

	voteResponse, err := raftServer2.RequestVote(ctx, request)

	if err != nil {
		t.Fatalf("RequestVote failed: %v", err)
	}

	expectedVoteGranted := true
	if voteResponse.VoteGranted != expectedVoteGranted {
		t.Errorf("Expected VoteGranted: %v, but got: %v", expectedVoteGranted, voteResponse.VoteGranted)
	}

	if raftNode2.votedFor != request.CandidateId {
		t.Errorf("Expected votedFor: %s, but got: %s", request.CandidateId, raftNode2.votedFor)
	}

	// Rejected request
	raftNode2.currentTerm = 2

	voteResponse, err = raftServer2.RequestVote(ctx, request)
	if err != nil {
		t.Fatalf("RequestVote failed: %v", err)
	}
	expectedVoteGranted = false
	if voteResponse.VoteGranted != expectedVoteGranted {
		t.Errorf("Expected VoteGranted: %v, but got: %v", expectedVoteGranted, voteResponse.VoteGranted)
	}
}

func startGRPCServer(node *RaftNode, lis net.Listener) {
	grpcServer := grpc.NewServer()
	pb.RegisterRaftServiceServer(grpcServer, &ServiceServer{Server: node})
	pb.RegisterRaftServer(grpcServer, &RaftServer{Server: node})
	grpcServer.Serve(lis)
}
