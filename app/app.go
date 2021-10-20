package app

import (
	"errors"
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	keepContainer, err :=  strconv.ParseBool(c.Query("keep"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keep value must be either true or false"})
		return
	}

	err = app.controller.UpdateContainer(containerName, image, keepContainer)
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

func (app *App) RollbackContainer(c *gin.Context) {
	containerName := c.Query("container")

	if containerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container must not be empty"})
		return
	}

	err := app.controller.RollbackContainer(containerName)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrContainerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "this container does not have a rollback container"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully restored container"})
}
