package handlers

import (
	"fmt"
	"github.com/XiovV/docker_control/controller"
	"github.com/XiovV/docker_control/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateHandler struct {
	controller *controller.DockerController
}

func NewUpdateHandler(controller *controller.DockerController) *UpdateHandler {
	return &UpdateHandler{controller: controller}
}

type UpdateRequest struct {
	Image     string `json:"image"`
	Container string `json:"container"`
}

type PullImageRequest struct {
	Image string `json:"image"`
}

func (uh *UpdateHandler) ContainerUpdate(c *gin.Context) {
	var updateRequest UpdateRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	containerId := uh.controller.FindContainerIDByName(updateRequest.Container)
	fmt.Println(containerId)

	if err := uh.controller.UpdateContainer(containerId, updateRequest.Image); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (uh *UpdateHandler) PullImage(c *gin.Context) {
	var pullImageRequest PullImageRequest

	if err := c.ShouldBindJSON(&pullImageRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := uh.controller.PullImage(pullImageRequest.Image); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	fmt.Println("pulled image:", pullImageRequest.Image)
}

func (uh *UpdateHandler) NodeStatus(c *gin.Context) {
	var nodeStatusRequest models.NodeStatusRequest

	if err := c.ShouldBindJSON(&nodeStatusRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	containerStatus, ok := uh.controller.GetContainerStatus(nodeStatusRequest.Container)

	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, containerStatus)
}

func (uh *UpdateHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.1.0"})
}