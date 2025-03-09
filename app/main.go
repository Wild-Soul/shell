package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func processCommand(cmd string) {
	if strings.HasPrefix(cmd, "exit") {
		os.Exit(0)

	} else if strings.HasPrefix(cmd, "echo") {

		elemap := mapFunction(strings.Split(cmd, "echo"), func(ele string) string {
			return strings.TrimLeft(ele, " ")
		})

		out := filterFunction(elemap, func(ele string) bool {
			return len(ele) != 0
		})

		fmt.Println(strings.Join(out, ""))

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
	processCommand(command)
}

func main() {
	// Start read-eval-print loop.
	for {
		getUserInput()
	}
}
