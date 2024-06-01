package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
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

func SetTargetServer(server *t.Address) {
	if conn != nil {
		conn.Close()
	}

	targetServer = server
	conn, err = grpc.NewClient(targetServer.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(t.LoggingInterceptorClient(ClientLogger)))
	if err != nil {
		return
	}

	serviceClient = pb.NewRaftServiceClient(conn)
}

func ExecuteCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	command := args[0]
	commandArgs := args[1:]

	if err := ValidateCommand(command, commandArgs); err != nil {
		fmt.Println(err)
		fmt.Println()
		return
	}

	switch command {
	case "ping":
		Ping()
	case "get":
		Get(commandArgs)
	case "set":
		Set(commandArgs)
	case "strln":
		Strln(commandArgs)
	case "del":
		Del(commandArgs)
	case "append":
		Append(commandArgs)
	case "request-log":
		RequestLog()
	case "change-server":
		ChangeServer(commandArgs)
	case "listen-web":
		ListenWeb(commandArgs)
	case "help":
		Help()
	case "exit":
		Exit()
	default:
	}

	fmt.Println()
}

func ValidateCommand(command string, args []string) error {
	switch command {
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
	default:
		return errors.New("unknown command: " + command)
	}

	return nil
}

func Ping() string {
	fmt.Println("Pinging", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	return r.GetResponse()
}

func Get(args []string) string {
	// TODO
	fmt.Println("Geting value", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Get(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	return r.GetValue()
}

func Set(args []string) string {
	// TODO
	fmt.Println("Set value", args[0], "with", args[1])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Set(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	return r.GetResponse()
}

func Strln(args []string) string {
	// TODO
	fmt.Println("Strln", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Strln(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	return r.GetValue()
}

func Del(args []string) string {
	// TODO
	fmt.Println("Delete value", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Del(ctx, &pb.KeyedRequest{Key: args[0]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetValue())
	return r.GetValue()
}

func Append(args []string) string {
	// TODO
	fmt.Println("Append value", args[0], "with", args[1])

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.Append(ctx, &pb.KeyValuedRequest{Key: args[0], Value: args[1]})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetResponse())
	return r.GetResponse()
}

func RequestLog() string {
	// TODO
	fmt.Println("Request Log", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT)
	defer cancel()

	r, err := serviceClient.ReqLog(ctx, &pb.LogRequest{})
	if err != nil {
		ClientLogger.Println(err)
		return err.Error()
	}

	fmt.Println(r.GetLogEntries())
	// TODO print log entries as messages
	return ""
}

func ChangeServer(args []string) string {
	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return err.Error()
	}

	newServer := t.NewAddress(host, port)
	SetTargetServer(&newServer)
	fmt.Println("Successfully changing server to", targetServer)
	return fmt.Sprintf("Successfully changing server to", targetServer)
}

func ListenWeb(args []string) {
	host := args[0]
	port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return
	}
	address := t.NewAddress(host, port)
	
	// TODO: correct HTTP method
	http.HandleFunc("/ping", LoggingMiddleware(handlePing))
	http.HandleFunc("/get", LoggingMiddleware(handleGet))
	http.HandleFunc("/set", LoggingMiddleware(handleSet))
	http.HandleFunc("/strln", LoggingMiddleware(handleStrln))
	http.HandleFunc("/del", LoggingMiddleware(handleDel))
	http.HandleFunc("/append", LoggingMiddleware(handleAppend))
	http.HandleFunc("/request-log", LoggingMiddleware(handleRequestLog))

	fmt.Println("Starting HTTP server on", address)
	err = http.ListenAndServe(address.String(), nil)
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func Help() {
	fmt.Println("Available commands:")
	fmt.Println("ping          : Ping the server.")
	fmt.Println("get           : Get the value associated with a key.")
	fmt.Println("set           : Set a key-value pair.")
	fmt.Println("strln         : Get the length of a string value associated with a key.")
	fmt.Println("del           : Delete a key-value pair.")
	fmt.Println("append        : Append a value to the string value associated with a key.")
	fmt.Println("request-log   : Request the log entries from the server.")
	fmt.Println("change-server : Change the target server address.")
	fmt.Println("listen-web    : Change client into web server.")
	fmt.Println("help          : Display available commands and their descriptions.")
	fmt.Println("exit          : Exit the program.")
}

func Exit() {
	fmt.Println("Exiting client...")
	os.Exit(0)
}
