package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var supportedCmds map[string]bool

func mapFunction[T any, U any](collection []T, f func(T) U) []U {
	var result []U

	for _, ele := range collection {
		result = append(result, f(ele))
	}

	return result
}

func filterFunction[T any](collection []T, f func(T) bool) []T {
	var result []T

	for _, ele := range collection {
		if f(ele) {
			result = append(result, ele)
		}
	}

	return result
}

func processCommand(cmd, args string) {
	if _, ok := supportedCmds[cmd]; ok {
		switch cmd {
		case "exit":
			os.Exit(0)

		case "echo":
			out := filterFunction(strings.Split(args, " "), func(ele string) bool {
				return len(ele) != 0
			})

			fmt.Println(strings.Join(out, " "))

		case "type":
			if _, present := supportedCmds[args]; present {
				fmt.Println(args + " is a shell builtin")
			} else {
				fmt.Println(args + ": not found")
			}
		}
	} else {
		fmt.Println(cmd + ": command not found")
	}
}

func getUserInput() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error while reading command: %v", err)
		os.Exit(1)
	}
	command = command[:len(command)-1]
	command = strings.Trim(command, " ")

	parts := strings.Fields(command)
	cmd := parts[0]
	args := strings.Join(parts[1:], " ")
	processCommand(cmd, args)
}

func initSupportedCmds() {
	supportedCmds = map[string]bool{
		"echo": true,
		"type": true,
		"exit": true,
	}
}

func main() {
	// Init supported commands
	initSupportedCmds()

	// Start read-eval-print loop.
	for {
		getUserInput()
	}
}
