package controller

import (
	"errors"
	"fmt"
)

var (
	ErrContainerNotRunning       = errors.New("container is not running")
	ErrContainerRestoreFailed    = errors.New("couldn't restore container")
	ErrContainerNotFound         = errors.New("container does not exist")
	ErrRollbackContainerNotFound = errors.New("rollback container does not exist ")
	ErrImageFormatInvalid        = errors.New("image format is invalid")
)

type ErrContainerStartFailed struct {
	ContainerId string
	Reason      error
}

func (e ErrContainerStartFailed) Error() string {
	return fmt.Sprintf("container %s could not be started: %s", e.ContainerId, e.Reason)
}
