package app

import (
	"github.com/XiovV/dokkup-agent/config"
	"github.com/XiovV/dokkup-agent/controller"
)

type App struct {
	controller controller.ContainerController
	config     *config.Config
}

func New(controller controller.ContainerController, config *config.Config) *App {
	return &App{controller: controller, config: config}
}
