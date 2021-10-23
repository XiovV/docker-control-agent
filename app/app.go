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

func (app *App) GetContainerImage(c *gin.Context) {
	containerName := c.Param("containerName")

	if containerName == "" {
		app.notFoundErrorResponse(c, "container not found")
		return
	}

	container, ok := app.controller.FindContainerByName(containerName)
	if !ok {
		app.notFoundErrorResponse(c, "container not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": container.Image})
}

func (app *App) PullImage(c *gin.Context) {
	image := c.Query("image")

	if image == "" {
		app.badRequestResponse(c, "image value must not be empty")
		return
	}

	err := app.controller.PullImage(image)
	if err != nil {
		switch {
		case errors.Is(err, controller.ErrImageFormatInvalid):
			app.badRequestResponse(c, "image format is invalid")
		default:
			app.internalErrorResponse(c, err.Error())
		}
		return
	}
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

	keepContainer, err := strconv.ParseBool(c.Query("keep"))
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
			app.notFoundErrorResponse(c, "the requested container does not exist")
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
