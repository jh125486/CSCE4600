package main

import (
	"bufio"
	"fmt"
	"github.com/jh125486/CSCE4600/Project2/builtins"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	var (
		input    string
		readLoop = bufio.NewReader(os.Stdin)
	)
	for {
		if err := printPrompt(os.Stdout, u); err != nil {
			panic(err)
		}
		if input, err = readLoop.ReadString('\n'); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		if err = handleInput(os.Stdout, input); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}
}

func printPrompt(w io.Writer, u *user.User) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "exit":
		os.Exit(0)
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
