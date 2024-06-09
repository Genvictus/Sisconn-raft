package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
	"bytes"
	"log"
	"net"
	"net/http"

	r "Sisconn-raft/raft"

	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var IsTransactionStart bool = false

type Transaction struct {
	Operation string
	Key       string
	Value     string
}

type NodeResponse struct {
	Address string
	State   string
}

var TransactionList []Transaction

func setClientLogger(serverAddress *t.Address) {
	ClientLogger = log.New(os.Stdout, "[Raft] Client "+serverAddress.String()+" : ", 0)
}

func setTargetServer(server *t.Address) {
	if conn != nil {
		conn.Close()
	}

	targetServer = server
	setClientLogger(targetServer)

	conn, err = grpc.NewClient(targetServer.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(t.LoggingInterceptorClient(ClientLogger)))
	if err != nil {
		return
	}

	serviceClient = pb.NewRaftServiceClient(conn)
}

func changeToLeaderAddress(leaderAddress string) error {
	host, port, err := net.SplitHostPort(leaderAddress)
	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	newServer := t.NewAddress(host, portInt)
	setTargetServer(&newServer)

	return nil
}

func executeCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	command := args[0]
	commandArgs := args[1:]

	if err := validateCommand(command, commandArgs); err != nil {
		fmt.Println(err)
		fmt.Println()
		return
	}

	switch command {
	case "start":
		start()
	case "ping":
		ping()
	case "get":
		get(commandArgs)
	case "set":
		set(commandArgs)
	case "strln":
		strln(commandArgs)
	case "del":
		del(commandArgs)
	case "append":
		Append(commandArgs)
	case "request-log":
		requestLog(serviceClient)
	case "change-server":
		changeServer(commandArgs)
	case "listen-web":
		listenWeb(commandArgs)
	case "help":
		help()
	case "exit":
		exit()
	case "commit":
		Commit()
	case "cancel":
		Cancel()
	case "add-node":
		AddNode(commandArgs)
	case "remove-node":
		RemoveNode(commandArgs)
	case "cluster-info":
		GetAllNode(commandArgs)
	default:
	}

	fmt.Println()
}

func validateCommand(command string, args []string) error {
	switch command {
	case "start":
		if IsTransactionStart {
			return errors.New("transaction already start")
		}
	case "ping":
		if len(args) != 0 {
			return errors.New("usage: ping")
		}
	case "get":
		if len(args) != 1 {
			return errors.New("usage: get <key>")
		}
	case "set":
		if len(args) != 2 {
			return errors.New("usage: set <key> <value>")
		}
	case "strln":
		if len(args) != 1 {
			return errors.New("usage: strln <key>")
		}
	case "del":
		if len(args) != 1 {
			return errors.New("usage: del <key>")
		}
	case "append":
		if len(args) != 2 {
			return errors.New("usage: append <key> <value>")
		}
	case "request-log":
		if len(args) != 0 {
			return errors.New("usage: request-log")
		}
	case "change-server":
		if len(args) != 2 {
			return errors.New("usage: change-server <host> <port>")
		}
	case "listen-web":
		if len(args) != 2 {
			return errors.New("usage: listen-web <host> <port>")
		}
	case "help":
		if len(args) != 0 {
			return errors.New("usage: help")
		}
	case "exit":
		if len(args) != 0 {
			return errors.New("usage: exit")
		}
	case "commit":
		if !IsTransactionStart {
			return errors.New("cant use commit")
		}
	case "cancel":
		if !IsTransactionStart {
			return errors.New("transaction is not started yet")
		}
	case "add-node":
		if len(args) != 2 {
			return errors.New("usage: add-node <host> <port>")
		}
	case "remove-node":
		if len(args) != 2 {
			return errors.New("usage: remove-node <host> <port>")
		}
	case "cluster-info":
		if len(args) != 0 {
			return errors.New("usage: cluster-info")
		}
	default:
		return errors.New("unknown command: " + command)
	}

	return nil
}

func start() {
	IsTransactionStart = true
	fmt.Println("Start Transaction Batch")
}

func ping() string {
	if IsTransactionStart {
		transaction := Transaction{
			Operation: "ping",
			Key:       "",
			Value:     "",
		}

		TransactionList = append(TransactionList, transaction)
		return ""
	}
	fmt.Println("Pinging", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Ping(ctx, &pb.PingRequest{})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
	}
	return r.GetResponse()
}

func get(args []string) string {
	fmt.Println("Geting value", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Get(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Get(ctx, &pb.KeyedRequest{Key: args[0]})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetValue())
	}
	return r.GetValue()
}

func set(args []string) string {
	if IsTransactionStart {
		transaction := Transaction{
			Operation: "set",
			Key:       args[0],
			Value:     args[1],
		}

		TransactionList = append(TransactionList, transaction)
		return ""
	}
	fmt.Println("Set value", args[0], "with", args[1])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Set(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Set(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetResponse())
	}

	return r.GetResponse()
}

func strln(args []string) string {
	fmt.Println("Strln", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Strln(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Strln(ctx, &pb.KeyedRequest{Key: args[0]})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetValue())
	}

	return r.GetValue()
}

func del(args []string) string {
	if IsTransactionStart {
		transaction := Transaction{
			Operation: "del",
			Key:       args[0],
			Value:     "",
		}

		TransactionList = append(TransactionList, transaction)
		return ""
	}

	fmt.Println("Delete value", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Del(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Del(ctx, &pb.KeyedRequest{Key: args[0]})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetValue())
	}

	return r.GetValue()
}

func Append(args []string) string {
	if IsTransactionStart {
		transaction := Transaction{
			Operation: "append",
			Key:       args[0],
			Value:     args[1],
		}

		TransactionList = append(TransactionList, transaction)
		return ""
	}

	fmt.Println("Append value", args[0], "with", args[1])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Append(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Append(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetResponse())
	}
	return r.GetResponse()
}

func Cancel() string {
	fmt.Println("transaction canceled")

	TransactionList = nil
	IsTransactionStart = false

	return ""
}

func Commit() string {
	fmt.Println("executing transaction batch")

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	var commitEntries []*pb.CommitEntry
	for _, transaction := range TransactionList {
		commitEntries = append(commitEntries, &pb.CommitEntry{
			Type:  transaction.Operation,
			Key:   transaction.Key,
			Value: transaction.Value,
		})
	}

	r, err := serviceClient.Commit(ctx, &pb.CommitRequest{
		CommitEntries: commitEntries,
	})

	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	// log.Println("buat ctx ga error", ctx)

	// log.Println("batch exe: ", TransactionList)
	fmt.Println(r.GetResponse())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.Commit(ctx, &pb.CommitRequest{
			CommitEntries: commitEntries,
		})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetResponse())
	}

	TransactionList = nil
	IsTransactionStart = false

	return r.GetResponse()
}

func AddNode(args []string) string {
	host := args[0]
	_, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return ""
	}
	fmt.Println("Adding new node")

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.AddNode(ctx, &pb.KeyValuedRequest{
		Key:   host + ":" + args[1],
		Value: "",
	})

	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err = serviceClient.AddNode(ctx, &pb.KeyValuedRequest{
			Key:   host + ":" + args[1],
			Value: "",
		})

		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetResponse())
	}
	return r.GetResponse()
}

func RemoveNode(args []string) string {
	fmt.Println("Removing node")
	host := args[0]
	_, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.RemoveNode(ctx, &pb.KeyedRequest{
		Key: host + ":" + args[1],
	})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}
	fmt.Println(r.GetResponse())

	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}

		r, err := serviceClient.RemoveNode(ctx, &pb.KeyedRequest{
			Key: host + ":" + args[1],
		})
		if err != nil {
			ClientLogger.Println(err)
			return err.Error()
		}
		fmt.Println(r.GetResponse())
	}
	return r.GetResponse()
}

func requestLog(client pb.RaftServiceClient) string {
	fmt.Println("Request Log", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := client.ReqLog(ctx, &pb.LogRequest{})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	var buffer bytes.Buffer
	buffer.WriteString("Log Entries:\n")
	for i, entry := range r.GetLogEntries() {
		buffer.WriteString(fmt.Sprintf("Entry %d:\n", i+1))
		buffer.WriteString(fmt.Sprintf("\tTerm: %d\n", entry.Term))
		buffer.WriteString(fmt.Sprintf("\tKey: %s\n", entry.Key))
		buffer.WriteString(fmt.Sprintf("\tValue: %s\n", entry.Value))
	}
	formattedLog := buffer.String()
	fmt.Println(formattedLog)

	return formattedLog
}

func GetAllNode(args []string) ([]NodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.ReqClusterInfo(ctx, &pb.ClusterInfoRequest{})
	if err != nil {
		ClientLogger.Println(err)
		return nil, err
	}

	if r.LeaderAddress != "" {
		err = changeToLeaderAddress(r.LeaderAddress)
		if err != nil {
			ClientLogger.Println(err)
			return nil, err
		}

		r, err = serviceClient.ReqClusterInfo(ctx, &pb.ClusterInfoRequest{})
		if err != nil {
			ClientLogger.Println(err)
			return nil, err
		}

		// for _, nodeInfo := range r.NodeInfo {
		// 	fmt.Println(nodeInfo.Address)
		// }
	}

	var buffer bytes.Buffer
	buffer.WriteString("Node:\n")
	// fmt.Println(len(r.NodeInfo))
	var responses []NodeResponse
	for _, node := range r.NodeInfo {
		err = changeToLeaderAddress(node.Address)
		if err != nil {
			ClientLogger.Println(err)
			return nil, err
		}

		stateResponse, _ := serviceClient.NodeState(ctx, &pb.StateRequest{})

		var state string
		if stateResponse == nil {
			state = "INACTIVE"
		} else {
			state = stateResponse.State
		}
		buffer.WriteString(fmt.Sprintf("%s: %s\n", node.Address, state))

		responses = append(responses, NodeResponse{
			Address: node.Address,
			State:   state,
		})
	}
	formattedNode := buffer.String()
	fmt.Println(formattedNode)

	return responses, nil
}

func changeServer(args []string) {
	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return
	}

	newServer := t.NewAddress(host, port)
	setTargetServer(&newServer)
	fmt.Println("Successfully changing server to", targetServer)
}

func listenWeb(args []string) {
	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return
	}
	address := t.NewAddress(host, port)

	// TODO: correct HTTP method
	//Client Web
	http.HandleFunc("/ping", loggingMiddleware(handlePing))
	http.HandleFunc("/get", loggingMiddleware(handleGet))
	http.HandleFunc("/set", loggingMiddleware(handleSet))
	http.HandleFunc("/strln", loggingMiddleware(handleStrln))
	http.HandleFunc("/del", loggingMiddleware(handleDel))
	http.HandleFunc("/append", loggingMiddleware(handleAppend))
	http.HandleFunc("/request-log", loggingMiddleware(handleRequestLog))

	// Management Dashboard
	http.HandleFunc("/node", loggingMiddleware(handleGetAllNode))
	http.HandleFunc("/log", loggingMiddleware(handleGetLog))
	http.HandleFunc("/node/delete", loggingMiddleware(handleRemoveNode))
	http.HandleFunc("/node/add", loggingMiddleware(handleAddNode))

	fmt.Println("Starting HTTP server on", &address)
	err = http.ListenAndServe(address.String(), nil)
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func help() {
	fmt.Println("Available commands:")
	fmt.Println("ping           : Ping the server.")
	fmt.Println("get            : Get the value associated with a key.")
	fmt.Println("set            : Set a key-value pair.")
	fmt.Println("strln          : Get the length of a string value associated with a key.")
	fmt.Println("del            : Delete a key-value pair.")
	fmt.Println("append         : Append a value to the string value associated with a key.")
	fmt.Println("request-log    : Request the log entries from the server.")
	fmt.Println("change-server  : Change the target server address.")
	fmt.Println("listen-web     : Change client into web server.")
	fmt.Println("help           : Display available commands and their descriptions.")
	fmt.Println("start          : Start batch transaction")
	fmt.Println("commit         : Commit and execute batch transaction")
	fmt.Println("exit           : Exit the program.")
	fmt.Println("add-node       : Add new raft node")
	fmt.Println("remove-node    : Remove raft node")
	fmt.Println("cluster-info   : Get all node address and state")
}

func exit() {
	fmt.Println("Exiting client...")
	os.Exit(0)
}
