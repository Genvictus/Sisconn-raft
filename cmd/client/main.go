package main

import (
	"Sisconn-raft/raft/transport"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	server_ip := args[0]
	server_port, err := strconv.Atoi(args[1])
    if err != nil {
		fmt.Println("Error: invalid server port, port must be an integer")
        panic(err)
    }

	server_address := transport.NewAddress(server_ip, server_port)

	fmt.Println("Client Started")
	fmt.Println("Connecting to server at", &server_address)
}
