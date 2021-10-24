package app

import (
	"errors"
	"github.com/XiovV/docker_control/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

	app.successResponse(c, "image pulled successfully")
}
