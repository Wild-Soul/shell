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
	"pwd":  true,
	"cd":   true,
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
	targetCmd := Command{
		cmd: c.args[0],
	}
	if targetCmd.isBuiltinCmd() {
		fmt.Println(targetCmd.cmd + " is a shell builtin")
	} else {
		for _, path := range binPaths {
			actualMap, _ := commandsInPaths.Load(path)
			// fmt.Println("Path:", path, actualMap)
			if actualMap != nil {
				mapForPath := actualMap.(map[string]bool)
				if _, present := mapForPath[targetCmd.cmd]; present {
					fmt.Printf("%v is %v/%v\n", targetCmd.cmd, path, targetCmd.cmd)
					return
				}
			}
		}
		fmt.Println(targetCmd.cmd + ": not found")
	}
}

func (c *Command) changeDirHandler() {
	targetPath := c.args[0]
	if len(targetPath) == 0 {
		targetPath = "~" // If users just does cd we should switch to $HOME path.
	} else {
		// check if targetpath starts with ~ . in which case prepend $HOME value to path.
		if strings.HasPrefix(targetPath, "~") {
			targetPath = os.Getenv("HOME") + strings.TrimPrefix(targetPath, "~")
		}
	}

	if err := os.Chdir(targetPath); err != nil {
		fmt.Printf("cd: %v: No such file or directory\n", targetPath)
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
	case "cd":
		c.changeDirHandler()
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
