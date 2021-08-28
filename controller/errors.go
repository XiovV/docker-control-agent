package controller

import "errors"

var (
	ErrContainerNotRunning = errors.New("container is not running")
	ErrRestoreFailed = errors.New("couldn't restore container")
)
