package app

import (
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
)

type App struct {
	controller controller.ContainerController
	config     *config.Config
}

func New(controller controller.ContainerController, config *config.Config) *App {
	return &App{controller: controller, config: config}
}
