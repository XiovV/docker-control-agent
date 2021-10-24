// Package controller is a wrapper for the docker API
package controller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
)

const (
	RollbackContainerSuffix = "-rollback"
)

type ContainerController interface {
	FindContainerByName(string) (types.Container, bool)
	FindContainerIDByName(string) (string, bool)
	PullImage(string) error
	UpdateContainer(string, string, bool) error
	RollbackContainer(string) error
}

// OldContainerConfig holds the configuration settings of a container
// that's being updated which will then be copied over to the updated container
type OldContainerConfig struct {
	ContainerName       string
	ContainerConfig     *container.Config
	ContainerHostConfig *container.HostConfig
}

type DockerController struct {
	cli *client.Client
	ctx context.Context
}

// New returns a pointer to DockerController.
// Will panic if a new docker client couldn't be established
func New() *DockerController {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &DockerController{cli: cli, ctx: context.Background()}
}

// FindContainerByName is used for finding a container by its name.
// Note: it only searches through running containers.
// If a container is found, it will return a types.Container and true, signifying that the
// container has been found
func (dc *DockerController) FindContainerByName(containerName string) (types.Container, bool) {
	containers, err := dc.cli.ContainerList(dc.ctx, types.ContainerListOptions{})
	if err != nil {
		return types.Container{}, false
	}

	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			return container, true
		}
	}

	return types.Container{}, false
}

// FindContainerIDByName tries to find a container by its name, and if it's found,
// it will return a string containing the name of the container and true, signifying that
// the container has been found. Unlike FindContainerByName, this method searches through
// both running and stopped containers
func (dc *DockerController) FindContainerIDByName(containerName string) (string, bool) {
	containers, err := dc.cli.ContainerList(dc.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", false
	}
	for _, container := range containers {
		if container.Names[0][1:] == containerName {
			return container.ID, true
		}
	}

	return "", false
}

// copyContainerConfig gets a copy of the config for a container with a specific id
func (dc *DockerController) copyContainerConfig(containerId string) (OldContainerConfig, error) {
	containerJson, err := dc.cli.ContainerInspect(dc.ctx, containerId)

	if err != nil {
		return OldContainerConfig{}, err
	}

	return OldContainerConfig{
		ContainerConfig:     containerJson.Config,
		ContainerHostConfig: containerJson.HostConfig,
		ContainerName:       containerJson.ContainerJSONBase.Name,
	}, nil
}

// PullImage pulls a requested image. It will return an ErrImageFormatInvalid
// if the image is not in this format: imagename:tag. It checks if the requested
// image already exists, and if it does it returns immediately.
func (dc *DockerController) PullImage(image string) error {
	imageParts := strings.Split(image, ":")

	if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
		return ErrImageFormatInvalid
	}

	if dc.doesImageExist(image) {
		return nil
	}

	reader, err := dc.cli.ImagePull(dc.ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	if _, err = io.Copy(os.Stdout, reader); err != nil {
		return err
	}

	return nil
}

// doesImageExist goes through all images and checks if the requested image exists.
func (dc *DockerController) doesImageExist(image string) bool {
	if strings.Split(image, ":")[1] == "latest" {
		return false
	}

	images, err := dc.cli.ImageList(dc.ctx, types.ImageListOptions{All: true})
	if err != nil {
		fmt.Println("error while fetching images:", err)
		return false
	}

	for _, foundImage := range images {
		if len(foundImage.RepoTags) > 0 {
			if foundImage.RepoTags[0] == image {
				return true
			}
		}
	}

	return false
}

// RollbackContainer tries to find a container with the '-rollback' suffix in its name.
// If it finds one, it will remove the '-rollback' suffix and run it, and it will
// remove the previous container. It will return ErrContainerNotFound if the requested container
// doesn't exist, and ErrRollbackContainerNotFound if the requested container doesn't have
// its own fallback container.
func (dc *DockerController) RollbackContainer(containerName string) error {
	rollbackContainerId, ok := dc.FindContainerIDByName(containerName + RollbackContainerSuffix)
	if !ok {
		return ErrContainerNotFound
	}

	currentContainerId, ok := dc.FindContainerIDByName(containerName)
	if !ok {
		return ErrRollbackContainerNotFound
	}

	if err := dc.stopContainer(currentContainerId); err != nil {
		return fmt.Errorf("couldn't stop container %s: %w", containerName, err)
	}

	err := dc.removeContainer(currentContainerId)
	if err != nil {
		return fmt.Errorf("couldn't remove container %s: %w", currentContainerId, err)
	}

	err = dc.renameContainer(rollbackContainerId, containerName)
	if err != nil {
		return fmt.Errorf("couldn't rename container %s: %w", rollbackContainerId, err)
	}

	err = dc.startContainer(rollbackContainerId)
	if err != nil {
		return ErrContainerStartFailed{ContainerId: rollbackContainerId, Reason: err}
	}

	if !dc.isContainerRunning(rollbackContainerId) {
		fmt.Println("rollback container is not running...")
		return ErrContainerNotRunning
	}

	return nil
}

func (dc *DockerController) UpdateContainer(containerName, image string, keepContainer bool) error {
	imageParts := strings.Split(image, ":")

	if len(imageParts) != 2 || imageParts[0] == "" || imageParts[1] == "" {
		return ErrImageFormatInvalid
	}

	containerId, ok := dc.FindContainerIDByName(containerName)
	if !ok {
		return ErrContainerNotFound
	}

	rollbackContainerId, ok := dc.FindContainerIDByName(containerName + RollbackContainerSuffix)
	if ok {
		fmt.Println("removing rollback container")
		err := dc.removeContainer(rollbackContainerId)
		if err != nil {
			return fmt.Errorf("could not remove rollback container: %w", err)
		}
	}
	fmt.Println("rollback container doesn't exist, continuing...")

	configCopy, err := dc.copyContainerConfig(containerId)
	if err != nil {
		return fmt.Errorf("couldn't copy container config: %w", err)
	}

	fmt.Printf("renaming %s (%s) to %s-rollback\n", configCopy.ContainerName, containerId, configCopy.ContainerName)
	if err = dc.renameContainer(containerId, configCopy.ContainerName+RollbackContainerSuffix); err != nil {
		return fmt.Errorf("couldn't rename container: %w", err)
	}

	fmt.Println("creating new container...")
	newContainerId, err := dc.createContainer(configCopy, image)
	if err != nil {
		fmt.Println("couldn't create new container:", err)
		if err = dc.restoreContainer(containerId, newContainerId, configCopy.ContainerName); err != nil {
			return fmt.Errorf("couldn't restore old container: %w", err)
		}
		return err
	}

	fmt.Println("updated container id:", newContainerId)

	fmt.Printf("stopping %s-rollback (%s)\n", configCopy.ContainerName, containerId)
	if err = dc.stopContainer(containerId); err != nil {
		return fmt.Errorf("coulnd't stop container %s: %w", configCopy.ContainerName, err)
	}

	fmt.Printf("starting new container (%s)\n", newContainerId)
	if err = dc.startContainer(newContainerId); err != nil {
		return err
	}

	if !dc.isContainerRunning(newContainerId) {
		fmt.Println("new container is not running, trying to restore old container...")
		if err = dc.restoreContainer(containerId, newContainerId, configCopy.ContainerName); err != nil {
			return ErrContainerRestoreFailed
		}

		return ErrContainerNotRunning
	}

	if !keepContainer {
		fmt.Printf("removing container %s-rollback (%s)\n", configCopy.ContainerName, containerId)
		err = dc.removeContainer(containerId)
		if err != nil {
			return fmt.Errorf("couldn't remove container %s-rollback: %w", containerId, err)
		}
	}

	return nil
}

func (dc *DockerController) doesContainerIDExist(containerId string) bool {
	containers, err := dc.cli.ContainerList(dc.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if container.ID == containerId {
			return true
		}
	}

	return false
}

func (dc *DockerController) restoreContainer(oldContainerId, newContainerId, originalName string) error {
	if dc.doesContainerIDExist(newContainerId) {
		fmt.Printf("RESTORE: removing newly created container %s\n", newContainerId)
		if err := dc.removeContainer(newContainerId); err != nil {
			return err
		}
	}

	fmt.Printf("RESTORE: renaming %s to %s\n", oldContainerId, originalName)
	if err := dc.renameContainer(oldContainerId, originalName); err != nil {
		return err
	}

	fmt.Printf("RESTORE: starting container %s\n", oldContainerId)
	if err := dc.startContainer(oldContainerId); err != nil {
		return err
	}

	return nil
}

func (dc *DockerController) removeContainer(containerId string) error {
	return dc.cli.ContainerRemove(dc.ctx, containerId, types.ContainerRemoveOptions{})
}

func (dc *DockerController) stopContainer(containerId string) error {
	return dc.cli.ContainerStop(dc.ctx, containerId, nil)
}

func (dc *DockerController) renameContainer(containerId, newName string) error {
	return dc.cli.ContainerRename(dc.ctx, containerId, newName)
}

func (dc *DockerController) createContainer(config OldContainerConfig, image string) (string, error) {
	config.ContainerConfig.Image = image

	resp, err := dc.cli.ContainerCreate(dc.ctx, config.ContainerConfig, config.ContainerHostConfig, nil, nil, config.ContainerName)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (dc *DockerController) startContainer(containerId string) error {
	if err := dc.cli.ContainerStart(dc.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (dc *DockerController) isContainerRunning(containerId string) bool {
	containers, err := dc.cli.ContainerList(dc.ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if container.ID == containerId {
			return true
		}
	}

	return false
}
