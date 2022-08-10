package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	if err = cmd.Start(); err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	out, _ := io.ReadAll(stdout)
	er, _ := io.ReadAll(stderr)

	fmt.Println(string(out))
	fmt.Println(string(er))
}
