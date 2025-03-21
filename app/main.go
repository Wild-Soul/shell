package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var supportedCmds map[string]bool
var binPaths []string

// can use map with mutex as well, since sync.Map is more optimized for concurrent reads and not writes.
var commandsInPaths sync.Map

func getUserInput() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error while reading command: %v", err)
		os.Exit(1)
	}
	command = command[:len(command)-1]

	command, args := parseCommand(command)
	// fmt.Println("Parse cmd:", command, args)

	cmd := Command{
		cmd:  command,
		args: args,
	}
	cmd.processCommand()
}

func initSupportedCmds() {
	supportedCmds = map[string]bool{
		"echo": true,
		"type": true,
		"exit": true,
	}
}

func initPaths() {
	binPaths = getUniqueStrings(strings.Split(os.Getenv("PATH"), ":"))
	// fmt.Println("Unique paths:", len(binPaths), strings.Join(binPaths, "\n"))
	var binPathWg sync.WaitGroup // to wait for all goroutines launched

	for _, path := range binPaths {
		binPathWg.Add(1)

		actualMap, _ := commandsInPaths.LoadOrStore(path, make(map[string]bool))
		mapForPath := actualMap.(map[string]bool) // type assert to normal map

		// Launch multiple go routines to read filepaths.
		go func(path string) {
			// fmt.Println("Starting path:", path)
			// st := time.Now()
			defer binPathWg.Done()
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// fmt.Println("Error in initPaths: ", err)
					return fmt.Errorf("error reading file: %w", err)
				}

				if !info.IsDir() {
					mapForPath[info.Name()] = true
				}
				return nil
			})

			// fmt.Println("Done with path:", path, time.Since(st))
		}(path)
	}

	binPathWg.Wait() // wait the main go routine for all go routines to complete.
}

func main() {
	var initWg sync.WaitGroup // init wait groups.

	// Init supported commands.
	initWg.Add(1)
	go func() {
		defer initWg.Done()
		initSupportedCmds()
	}()

	// Init bin paths.
	initWg.Add(1)
	go func() {
		defer initWg.Done()
		initPaths()
	}()

	// wait for all init to happen
	initWg.Wait()

	// fmt.Println("binPaths", binPaths)
	// fmt.Println("commandsInPaths:")
	// commandsInPaths.Range(func(key, value any) bool {
	// 	// Print each key-value pair
	// 	actualValue := value.(map[string]bool)
	// 	fmt.Printf("%v: %v\n", key, len(actualValue))
	// 	return true
	// })

	// fmt.Println("Shell initialized!")
	// Start read-eval-print loop.
	for {
		getUserInput()
	}
}
