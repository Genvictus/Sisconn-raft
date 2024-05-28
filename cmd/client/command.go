package main

import (
	"errors"
	"fmt"
	"strings"
)

func ExecuteCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	command := args[0]
	commandArgs := args[1:]

	if err := validateCommand(command, commandArgs); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
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
	case "help":
		help()
	case "exit":
		exit()
	default:
	}
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

func help() {
	// TODO
}

func exit() {
	// TODO
}
