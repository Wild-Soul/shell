package main

import (
	"bufio"
	"fmt"
	"os"
)

func getUserInput() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error while reading command: %v", err)
		os.Exit(1)
	}

	fmt.Println(command[:len(command)-1] + ": command not found")
}

func main() {
	// Start read-eval-print loop.
	for {
		getUserInput()
	}
}
