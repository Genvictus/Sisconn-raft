package main

import (
	t "Sisconn-raft/raft/transport"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ExecuteCommand(server *t.Address, input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	command := args[0]
	commandArgs := args[1:]

	if err := validateCommand(command, commandArgs); err != nil {
		fmt.Println(err)
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
	case "change-server":
		changeServer(commandArgs, server)
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
	// TODO
}

func get(args []string) {
	// TODO
}

func set(args []string) {
	// TODO
}

func strln(args []string) {
	// TODO
}

func del(args []string) {
	// TODO
}

func append(args []string) {
	// TODO
}

func changeServer(args []string, server *t.Address) {
	host := args[0]
	port, err := strconv.Atoi(args[1])
    if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
		return
    }

	server.IP = host
	server.Port = port
	fmt.Println("Successfully changing server to", server)
}

func help() {
	fmt.Println("Available commands:")
	fmt.Println("ping          : Ping the server.")
	fmt.Println("get           : Get the value associated with a key.")
	fmt.Println("set           : Set a key-value pair.")
	fmt.Println("strln         : Get the length of a string value associated with a key.")
	fmt.Println("del           : Delete a key-value pair.")
	fmt.Println("append        : Append a value to the string value associated with a key.")
	fmt.Println("change-server : Change the server address.")
	fmt.Println("help          : Display available commands and their descriptions.")
	fmt.Println("exit          : Exit the program.")
}

func exit() {
	fmt.Println("Exiting client...")
	os.Exit(0);
}
