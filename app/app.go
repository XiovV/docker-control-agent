package app

import (
	"fmt"
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
	"github.com/XiovV/docker_control/models"
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

type UpdateRequest struct {
	Image     string `json:"image"`
	Container string `json:"container"`
}

type PullImageRequest struct {
	Image string `json:"image"`
}

func (app *App) PullImage(c *gin.Context) {
	apiKey := c.GetHeader("key")
	fmt.Println(apiKey)

	image := c.Query("image")
	fmt.Println(image)
}

func (app *App) ContainerUpdate(c *gin.Context) {
	var updateRequest UpdateRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	containerId, ok := app.controller.FindContainerIDByName(updateRequest.Container)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	if err := app.controller.UpdateContainer(containerId, updateRequest.Image); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (app *App) NodeStatus(c *gin.Context) {
	var nodeStatusRequest models.NodeStatusRequest

	if err := c.ShouldBindJSON(&nodeStatusRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	containerStatus, ok := app.controller.GetContainerStatus(nodeStatusRequest.Container)

	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, containerStatus)
}

func (app *App) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.1.0"})
}