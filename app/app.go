package app

import (
	"errors"
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	controller *controller.DockerController
	config config.Config
}

func New(controller *controller.DockerController, config config.Config) *App {
	return &App{controller: controller, config: config}
}

func (app *App) PullImage(c *gin.Context) {
	image := c.Query("image")

	err := app.controller.PullImage(image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (app *App) UpdateContainer(c *gin.Context) {
	containerName := c.Query("container")
	image := c.Query("image")

	err := app.controller.UpdateContainer(containerName, image)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrContainerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "the requested resource could not be found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

}
