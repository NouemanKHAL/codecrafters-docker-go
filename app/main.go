package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"codecrafters-docker-go/app/api"
	"codecrafters-docker-go/app/util"

	"github.com/google/uuid"
)

const docker_explorer_path = "/usr/local/bin/docker-explorer"

func createContainer(path, image, tag string) error {
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

	d := api.NewDockerRegistryClient()
	manifest, err := d.PullImage(image, tag)
	if err != nil {
		return err
	}

	// download layers
	val := make(map[string]interface{})
	err = json.Unmarshal(manifest, &val)
	if err != nil {
		log.Fatal(err)
	}

	fsLayers := val["fsLayers"]
	repoName := val["name"].(string)

	for i, layer := range fsLayers.([]interface{}) {
		digest := layer.(map[string]interface{})["blobSum"].(string)
		layerBytes, err := d.PullLayer(repoName, digest)
		if err != nil {
			return err
		}
		digestPath := util.JoinPath(container_path, fmt.Sprintf("digest_%d.tar.gz", i+1))
		f, err := util.CreateFile(digestPath)
		if err != nil {
			return err
		}

		f.Write(layerBytes)

		tarCmd := exec.Command("tar", "-xvf", digestPath, "-C", container_path)
		if err := tarCmd.Run(); err != nil {
			return err
		}

		err = os.Remove(digestPath)
		if err != nil {
			return err
		}
	}

	err = syscall.Chroot(container_path)
	if err != nil {
		return fmt.Errorf("chroot: %s", err)
	}
	_, err = util.CreateFile("/dev/null")
	if err != nil {
		return fmt.Errorf("util.CreateFile: %s", err)
	}

	return err
}

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	imageArg := os.Args[2]
	var image, tag string
	if strings.Contains(imageArg, ":") {
		parts := strings.Split(os.Args[2], ":")
		image = parts[0]
		tag = parts[1]
	} else {
		image = imageArg
		tag = "latest"
	}

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	err := createContainer(".", image, tag)
	if err != nil {
		fmt.Printf("error creating a container: %v", err)
		os.Exit(1)
	}

	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		os.Exit(e.ExitCode())
	} else if err != nil {
		fmt.Printf("error while running command '%s': %v", cmd.String(), err)
		os.Exit(1)
	}
}
