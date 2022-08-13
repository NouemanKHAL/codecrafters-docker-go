package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"

	"codecrafters-docker-go/app/util"

	"github.com/google/uuid"
)

const docker_explorer_path = "/usr/local/bin/docker-explorer"

func createContainer(path string) error {
	container_id := uuid.New()
	container_path := util.JoinPath(path, container_id.String())

	err := os.Mkdir(container_path, 0750)
	if err != nil {
		return err
	}

	container_docker_explorer := util.JoinPath(container_path, docker_explorer_path)
	err = util.CopyFile(docker_explorer_path, container_docker_explorer)
	if err != nil {
		return err
	}

	cmd := exec.Cmd{
		Path:   "/bin/sh",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
		},
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = syscall.Chroot(container_path)
	if err != nil {
		return nil
	}
	_, err = util.CreateFile("/dev/null")
	return err
}

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	rand.Seed(time.Now().UnixNano())

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := createContainer(".")
	if err != nil {
		fmt.Printf("error creating a container: %v", err)
		os.Exit(1)
	}

	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		os.Exit(e.ExitCode())
	} else if err != nil {
		fmt.Printf("error while running command '%s': %v", cmd.String(), err)
		os.Exit(1)
	}
}
