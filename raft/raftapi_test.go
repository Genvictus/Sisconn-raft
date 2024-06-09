package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"context"
	"fmt"
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
	raftNode.initiateLeader()
	go raftNode.runTest()

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
			{Type: "del", Key: "CommitKey5", Value: ""},
			{Type: "del", Key: "CommitKey2", Value: ""},
			{Type: "append", Key: "CommitKey3", Value: " AppendValue3"},
			{Type: "set", Key: "CommitKey4", Value: "ResetValue4"},
			{Type: "append", Key: "CommitKey5", Value: " AppendValue5"},
			{Type: "append", Key: "CommitKey5", Value: " AppendValue5"},
		},
	}

	// Commit
	response, err := client.Commit(ctx, request)
	if err != nil {
		t.Errorf("Commit failed: %v", err)
	}
	expected := fmt.Sprintf("OK (%d commands execution)", len(request.CommitEntries))
	if response.Response != expected {
		t.Errorf("Expected response: %s, but got: %s", expected, response.Response)
	}

	// Check get value
	tests := []*pb.KeyValuedRequest{
		{Key: "CommitKey1", Value: "CommitValue1"},
		{Key: "CommitKey2", Value: ""},
		{Key: "CommitKey3", Value: "CommitValue3 AppendValue3"},
		{Key: "CommitKey4", Value: "ResetValue4"},
		{Key: "CommitKey5", Value: " AppendValue5 AppendValue5"},
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

	// Append nil value
	delRequest := &pb.KeyedRequest{Key: key}
	_, err = client.Del(ctx, delRequest)
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	appendResponse, err = client.Append(ctx, appendRequest)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	if appendResponse.Response != expectedResponse {
		t.Errorf("Expected response: %s, but got: %s", expectedResponse, appendResponse.Response)
	}

	getRequest = &pb.KeyedRequest{Key: key}
	getResponse, err = client.Get(ctx, getRequest)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	expectedValue = appendValue
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

	go raftNode1.runTest()

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
	go raftNode2.runTest()

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

func TestAppendEntries(t *testing.T) {
	// Setup server
	serverAddress := transport.NewAddress("localhost", 3011)
	lis, err := net.Listen("tcp", serverAddress.String())
	if err != nil {
		t.Fatalf("Failed to listen for the first server: %v", err)
	}
	defer lis.Close()

	followerAddress := transport.NewAddress("localhost", 3012)
	lis2, err := net.Listen("tcp", followerAddress.String())
	if err != nil {
		t.Fatalf("Failed to listen for the second server: %v", err)
	}
	defer lis2.Close()

	raftNode := NewNode(serverAddress.String())
	followerNode := NewNode(followerAddress.String())

	raftNode.AddConnections([]string{serverAddress.String(), followerAddress.String()})
	go startGRPCServer(raftNode, lis)
	go startGRPCServer(followerNode, lis2)

	raftServer := raftNode.raftClient[followerAddress.String()]
	go raftNode.runTest()
	go followerNode.runTest()

	// AppendEntries request
	request := &pb.AppendEntriesArg{
		Term:         1,
		LeaderId:     "leader1",
		PrevLogIndex: 0,
		PrevLogTerm:  0,

		LogEntries: []*pb.LogEntry{
			{Term: 1, Key: "key1", Value: "value1"},
			{Term: 2, Key: "key2", Value: "value2"},
		},

		LeaderCommit: 0,
		LogType: _DataLog,
	}

	ctx := context.Background()

	// Test apppend entries
	appendEntriesResponse, err := raftServer.AppendEntries(ctx, request)

	if err != nil {
		t.Fatalf("AppendEntries failed: %v", err)
	}

	expectedSuccess := true
	if appendEntriesResponse.Success != expectedSuccess {
		t.Errorf("Expected Success: %v, but got: %v", expectedSuccess, appendEntriesResponse.Success)
	}

	// Test outdated term
	followerNode.currentTerm = 69
	appendEntriesResponse, err = raftServer.AppendEntries(ctx, request)

	if err != nil {
		t.Fatalf("AppendEntries failed: %v", err)
	}

	expectedSuccess = false
	if appendEntriesResponse.Success != expectedSuccess {
		t.Errorf("Expected Success: %v, but got: %v", expectedSuccess, appendEntriesResponse.Success)
	}

	followerNode.currentTerm = 1

	// Test index out of bound
	followerNode.log.lastIndex = 69
	request.PrevLogIndex = 420

	appendEntriesResponse, err = raftServer.AppendEntries(ctx, request)

	if err != nil {
		t.Fatalf("AppendEntries failed: %v", err)
	}

	if appendEntriesResponse.Success != expectedSuccess {
		t.Errorf("Expected Success: %v, but got: %v", expectedSuccess, appendEntriesResponse.Success)
	}

	followerNode.log.lastIndex = 0
	request.PrevLogIndex = 0

	// Test invalid entry term
	request.PrevLogTerm = 69

	appendEntriesResponse, err = raftServer.AppendEntries(ctx, request)

	if err != nil {
		t.Fatalf("AppendEntries failed: %v", err)
	}

	if appendEntriesResponse.Success != expectedSuccess {
		t.Errorf("Expected Success: %v, but got: %v", expectedSuccess, appendEntriesResponse.Success)
	}

	request.PrevLogTerm = 0

	// Test invalid entry term
	request.LeaderCommit = 69

	appendEntriesResponse, err = raftServer.AppendEntries(ctx, request)

	if err != nil {
		t.Fatalf("AppendEntries failed: %v", err)
	}

	expectedSuccess = true

	if appendEntriesResponse.Success != expectedSuccess {
		t.Errorf("Expected Success: %v, but got: %v", expectedSuccess, appendEntriesResponse.Success)
	}
}

func TestFollowerRequestRedirect(t *testing.T) {
	// Setup follower
	leaderAddress := transport.NewAddress("localhost", 8000)
	followerAddress := transport.NewAddress("localhost", 8001)
	lis, err := net.Listen("tcp", followerAddress.String())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	followerNode := NewNode(followerAddress.String())
	followerNode.leaderAddress = leaderAddress.String()

	go startGRPCServer(followerNode, lis)

	conn, err := grpc.NewClient(followerAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	follower := pb.NewRaftServiceClient(conn)

	// Setup test
	ctx := context.Background()

	var (
		messageResponse *pb.MessageResponse
		valueResponse   *pb.ValueResponse
		logresponse     *pb.LogResponse
	)

	expected := NotLeaderResponse + leaderAddress.String()

	// Ping
	messageResponse, err = follower.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	if messageResponse.Response != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, messageResponse.Response)
	}

	if messageResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, messageResponse.LeaderAddress)
	}

	// Get
	valueResponse, err = follower.Get(ctx, &pb.KeyedRequest{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if valueResponse.Value != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, valueResponse.Value)
	}

	if valueResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, valueResponse.LeaderAddress)
	}

	// Set
	messageResponse, err = follower.Set(ctx, &pb.KeyValuedRequest{})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if messageResponse.Response != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, messageResponse.Response)
	}

	if messageResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, messageResponse.LeaderAddress)
	}

	// Strln
	valueResponse, err = follower.Strln(ctx, &pb.KeyedRequest{})
	if err != nil {
		t.Fatalf("Strln failed: %v", err)
	}

	if valueResponse.Value != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, valueResponse.Value)
	}

	if valueResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, valueResponse.LeaderAddress)
	}

	// Del
	valueResponse, err = follower.Del(ctx, &pb.KeyedRequest{})
	if err != nil {
		t.Fatalf("Del failed: %v", err)
	}

	if valueResponse.Value != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, valueResponse.Value)
	}

	if valueResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, valueResponse.LeaderAddress)
	}

	// Append
	messageResponse, err = follower.Append(ctx, &pb.KeyValuedRequest{})
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if messageResponse.Response != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, messageResponse.Response)
	}

	if messageResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, messageResponse.LeaderAddress)
	}

	// Append
	messageResponse, err = follower.Append(ctx, &pb.KeyValuedRequest{})
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if messageResponse.Response != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, messageResponse.Response)
	}

	if messageResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, messageResponse.LeaderAddress)
	}

	// ReqLog
	logresponse, err = follower.ReqLog(ctx, &pb.LogRequest{})
	if err != nil {
		t.Fatalf("ReqLog failed: %v", err)
	}

	if logresponse.LogEntries != nil {
		t.Errorf("Expected response: %v, but got: %v", nil, logresponse.LogEntries)
	}

	if logresponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, logresponse.LeaderAddress)
	}

	// Commit
	messageResponse, err = follower.Commit(ctx, &pb.CommitRequest{})
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	if messageResponse.Response != expected {
		t.Errorf("Expected response: %s, but got: %v", expected, messageResponse.Response)
	}

	if messageResponse.LeaderAddress != leaderAddress.String() {
		t.Errorf("Expected leader address: %s, but got: %s", &leaderAddress, messageResponse.LeaderAddress)
	}
}

func TestServerPendingRequest(t *testing.T) {
	// Setup server cluster
	serverAddress := transport.NewAddress("localhost", 8021)
	lis, err := net.Listen("tcp", serverAddress.String())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	raftNode := NewNode(serverAddress.String())
	raftNode.AddConnections([]string{serverAddress.String(), "localhost:8022"})
	raftNode.initiateLeader()
	go raftNode.runTest()

	go startGRPCServer(raftNode, lis)

	conn, err := grpc.NewClient(serverAddress.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	pendingClient := pb.NewRaftServiceClient(conn)

	// Setup test
	ctx := context.Background()

	// Set
	setRequest := &pb.KeyValuedRequest{Key: "key", Value: "value"}
	setResponse, err := pendingClient.Set(ctx, setRequest)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if setResponse.Response != "PENDING" {
		t.Errorf("Expected response: PENDING, but got: %v", setResponse.Response)
	}

	// Append
	appendRequest := &pb.KeyValuedRequest{Key: "key", Value: "value"}
	appendResponse, err := pendingClient.Append(ctx, appendRequest)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if appendResponse.Response != "PENDING" {
		t.Errorf("Expected response: PENDING, but got: %v", appendResponse.Response)
	}

	// Commit
	commitRequest := &pb.CommitRequest{
		CommitEntries: []*pb.CommitEntry{
			{Type: "set", Key: "CommitKey1", Value: "CommitValue1"},
			{Type: "set", Key: "CommitKey2", Value: "CommitValue2"},
			{Type: "set", Key: "CommitKey3", Value: "CommitValue3"},
		},
	}

	commitResponse, err := pendingClient.Commit(ctx, commitRequest)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	if commitResponse.Response != "PENDING (3 commands execution)" {
		t.Errorf("Expected response: PENDING (3 commands execution), but got: %v", commitResponse.Response)
	}
}

func startGRPCServer(node *RaftNode, lis net.Listener) {
	grpcServer := grpc.NewServer()
	pb.RegisterRaftServiceServer(grpcServer, &ServiceServer{Server: node})
	pb.RegisterRaftServer(grpcServer, &RaftServer{Server: node})
	grpcServer.Serve(lis)
}
