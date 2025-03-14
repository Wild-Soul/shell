package main

import (
	"fmt"
	"os"
	"strings"
)

// map of command and it's handler.
func exitHandler(args []string) {
	os.Exit(0)
}

func echoHandler(args []string) {
	out := filterFunction(args, func(ele string) bool {
		return len(ele) != 0
	})

	fmt.Println(strings.Join(out, " "))
}

func typeHandler(args []string) {
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

func getCmdHandler(cmd string) func([]string) {
	return map[string]func([]string){
		"exit": exitHandler,
		"echo": echoHandler,
		"type": typeHandler,
	}[cmd]
}
