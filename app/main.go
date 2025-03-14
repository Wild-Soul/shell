package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var supportedCmds map[string]bool
var binPaths []string

// can use map with mutex as well, since sync.Map is more optimized for concurrent reads and not writes.
var commandsInPaths sync.Map

func processCommand(cmd string, args []string) {
	if _, ok := supportedCmds[cmd]; ok {

		switch cmd {
		case "exit":
			os.Exit(0)

		case "echo":
			out := filterFunction(args, func(ele string) bool {
				return len(ele) != 0
			})

			fmt.Println(strings.Join(out, " "))

		case "type":
			targetCmd := args[0]
			if _, present := supportedCmds[targetCmd]; present {
				fmt.Println(targetCmd + " is a shell builtin")
			} else {
				for _, path := range binPaths {
					actualMap, _ := commandsInPaths.Load(path)
					// fmt.Println("Path:", path, actualMap)
					if actualMap != nil {
						mapForPath := actualMap.(map[string]bool)
						if _, present := mapForPath[targetCmd]; present {
							fmt.Printf("%v is %v/%v\n", targetCmd, path, targetCmd)
							return
						}
					}
				}
				fmt.Println(targetCmd + ": not found")
			}
		}
	} else if out, err := exec.Command(cmd, args...).Output(); err == nil {
		fmt.Println(string(out))
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
	processCommand(parts[0], parts[1:])
}

func initSupportedCmds() {
	supportedCmds = map[string]bool{
		"echo": true,
		"type": true,
		"exit": true,
	}
}

func initPaths() {
	binPaths = strings.Split(os.Getenv("PATH"), ":")
	var binPathWg sync.WaitGroup // to wait for all goroutines launched

	for _, path := range binPaths {
		binPathWg.Add(1)

		actualMap, _ := commandsInPaths.LoadOrStore(path, make(map[string]bool))
		mapForPath := actualMap.(map[string]bool) // type assert to normal map

		// Launch multiple go routines to read filepaths.
		go func(path string) {
			defer binPathWg.Done()
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					// fmt.Println("Error in initPaths: ", err)
					// return fmt.Errorf("error reading file: %w", err)
					return nil
				}

				// fmt.Println("Path:", path, " item:", info.Name())
				if !info.IsDir() {
					mapForPath[info.Name()] = true
				}
				return nil
			})
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

	// Start read-eval-print loop.
	for {
		getUserInput()
	}
}
