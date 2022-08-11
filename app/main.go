package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	container_dir := "mycontainer"
	err := os.Mkdir("mycontainer", 0750)
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	mkdirCommand := exec.Command("mkdir", "-p", container_dir+"/usr/local/bin")
	err = mkdirCommand.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	src := "/usr/local/bin/docker-explorer"
	dst := container_dir + "/usr/local/bin/"

	copyCommand := exec.Command("cp", src, dst)
	err = copyCommand.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	err = syscall.Chroot(container_dir)
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	os.MkdirAll("/dev", 0750)
	os.Open("/dev/null")
	os.Chdir("/dev")
	os.Create("null")

	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		os.Exit(e.ExitCode())
	} else if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}
}
