package main

import (
	"fmt"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	fmt.Println("Implement your program here")

	// Comment this section out to pass the first stage!
	//
	// command := os.Args[3]
	// args := os.Args[4:len(os.Args)]
	//
	// cmd := exec.Command(command, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	//
	// if err := cmd.Run(); err != nil {
	// 	if _, ok := err.(*exec.ExitError); ok {
	// 		os.Exit(cmd.ProcessState.ExitCode())
	// 	}
	// 	fmt.Printf("Err: %v", err)
	// 	os.Exit(1)
	// }
}
