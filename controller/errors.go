package controller

import "errors"

var (
	ErrContainerNotRunning = errors.New("container is not running")
	ErrContainerRestoreFailed = errors.New("couldn't restore container")
	ErrContainerNotFound = errors.New("couldn't find container")
)
