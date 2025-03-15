package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtInCmds = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
}

// Refactor idea: Create an interface for commands, with methods like Execute and Validate.
// This interface will be implemented by different Commands like EchoCommand, TypeCommand ...

type Command struct {
	cmd  string
	args []string
}

func (c *Command) exitHandler() {
	os.Exit(0)
}

func (c *Command) echoHandler() {
	out := filterFunction(c.args, func(ele string) bool {
		return len(ele) != 0
	})

	fmt.Println(strings.Join(out, " "))
}

func (c *Command) typeHandler() {
	targetCmd := c.args[0]
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

func (c *Command) pwdHandler() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error in pwd execution: %w", err)
	}

	fmt.Println(cwd)
	return nil
}

func (c *Command) isBuiltinCmd() bool {
	return builtInCmds[c.cmd]
}

func (c *Command) executeBuiltInCmd() {
	switch c.cmd {
	case "type":
		c.typeHandler()
	case "echo":
		c.echoHandler()
	case "exit":
		c.exitHandler()
	case "pwd":
		c.pwdHandler()
	}
}

func (c *Command) processCommand() {
	isBuiltIn := c.isBuiltinCmd()

	if isBuiltIn {
		c.executeBuiltInCmd()
	} else { // 2 possibilities: either the cmd is present and can be executed, or not found.
		for _, path := range binPaths {
			actualMap, _ := commandsInPaths.Load(path)
			commandsInPath := actualMap.(map[string]bool)
			if _, ok := commandsInPath[c.cmd]; ok {
				out, err := exec.Command(c.cmd, c.args...).Output()
				if err == nil {
					fmt.Print(string(out))
					return // found command stop looking
				} else {
					fmt.Println("ERROR", err.Error())
				}
			}
		}

		// Command not found at this point.
		fmt.Println(c.cmd + ": command not found")
	}
}
