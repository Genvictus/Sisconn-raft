package main

import (
	pb "Sisconn-raft/raft/raftpc"
	t "Sisconn-raft/raft/transport"
	"time"

	r "Sisconn-raft/raft"

	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setTargetServer(server *t.Address) {
	if conn != nil {
		conn.Close()
	}

	targetServer = server
	conn, err = grpc.NewClient(targetServer.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	if err := validateCommand(command, commandArgs); err != nil {
		fmt.Println(err)
		fmt.Println()
		return
	}

	switch command {
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
		append(commandArgs)
	case "request-log":
		requestLog()
	case "change-server":
		changeServer(commandArgs)
	case "help":
		help()
	case "exit":
		exit()
	default:
	}

	fmt.Println()
}

func validateCommand(command string, args []string) error {
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

func ping() {
	fmt.Println("Pinging", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetResponse())
}

func get(args []string) {
	// TODO
	fmt.Println("Geting value", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Get(ctx, &pb.KeyedRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetValue())
}

func set(args []string) {
	// TODO
	fmt.Println("Set value", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Set(ctx, &pb.KeyValuedRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetResponse())
}

func strln(args []string) {
	// TODO
	fmt.Println("Strln", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Strln(ctx, &pb.KeyedRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetValue())
}

func del(args []string) {
	// TODO
	fmt.Println("Set value", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Del(ctx, &pb.KeyedRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetValue())
}

func append(args []string) {
	// TODO
	fmt.Println("Append value", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.Append(ctx, &pb.KeyValuedRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetResponse())
}

func requestLog() {
	// TODO
	fmt.Println("Req Log", targetServer)

	ctx, cancel := context.WithTimeout(context.Background(), r.CLIENT_TIMEOUT*time.Millisecond)
	defer cancel()

	r, err := serviceClient.ReqLog(ctx, &pb.LogRequest{})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(r.GetLogEntries())
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

func help() {
	fmt.Println("Available commands:")
	fmt.Println("ping          : Ping the server.")
	fmt.Println("get           : Get the value associated with a key.")
	fmt.Println("set           : Set a key-value pair.")
	fmt.Println("strln         : Get the length of a string value associated with a key.")
	fmt.Println("del           : Delete a key-value pair.")
	fmt.Println("append        : Append a value to the string value associated with a key.")
	fmt.Println("request-log   : Change the server address.")
	fmt.Println("change-server : Change the server address.")
	fmt.Println("help          : Display available commands and their descriptions.")
	fmt.Println("exit          : Exit the program.")
}

func exit() {
	fmt.Println("Exiting client...")
	os.Exit(0)
}
