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

	if image == "" {
		app.badRequestResponse(c, "image value must not be empty")
		return
	}

	err := app.controller.PullImage(image)
	if err != nil {
		app.internalErrorResponse(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (app *App) UpdateContainer(c *gin.Context) {
	containerName := c.Query("container")
	image := c.Query("image")

	if containerName == "" {
		app.badRequestResponse(c, "container value must not be empty")
		return
	}

	if image == "" {
		app.badRequestResponse(c, "image value must not be empty")
		return
	}

	keepContainer, err :=  strconv.ParseBool(c.Query("keep"))
	if err != nil {
		app.badRequestResponse(c, "keep value must be either true or false")
		return
	}

	err = app.controller.UpdateContainer(containerName, image, keepContainer)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrContainerNotFound):
			app.notFoundErrorResponse(c, "the requested container could not be found")
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
}

func (app *App) RollbackContainer(c *gin.Context) {
	containerName := c.Query("container")

	if containerName == "" {
		app.badRequestResponse(c, "container value must not be empty")
		return
	}

	err := app.controller.RollbackContainer(containerName)
	if err != nil {
		var containerStartFailedErr controller.ErrContainerStartFailed

		switch {
		case errors.Is(err, controller.ErrContainerNotFound):
			app.notFoundErrorResponse(c, "the requested container could does not exist")
		case errors.Is(err, controller.ErrRollbackContainerNotFound):
			app.notFoundErrorResponse(c, "the requested container does not have a rollback container")
		case errors.Is(err, controller.ErrContainerNotRunning):
			app.internalErrorResponse(c, "the container failed to start")
		case errors.As(err, &containerStartFailedErr):
			app.internalErrorResponse(c, containerStartFailedErr.Reason.Error())
		default:
			app.internalErrorResponse(c, err.Error())
		}
		return
	}

	app.successResponse(c, "successfully restored container")
}
