package handlers

import (
	"fmt"
	"github.com/XiovV/docker_control/controller"
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
	ImageTag string `json:"image_tag"`
	Container string `json:"container"`
}

func (uh *UpdateHandler) HandleContainerUpdate(c *gin.Context) {
	var updateRequest UpdateRequest

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	containerId := uh.controller.FindContainerIDByName(updateRequest.Container)
	fmt.Println(containerId)

	if err := uh.controller.UpdateContainer(containerId, updateRequest.ImageTag); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}


}