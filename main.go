package main

import (
	"context"
	"github.com/XiovV/docker_control/controller"
	"github.com/XiovV/docker_control/handlers"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	dockerController := controller.New(cli, ctx)

	updateHandler := handlers.NewUpdateHandler(dockerController)

	router := gin.Default()
	router.POST("/api/containers/update", updateHandler.ContainerUpdate)
	router.POST("/api/images/pull", updateHandler.PullImage)
	router.POST("/api/nodes/status", updateHandler.NodeStatus)

	router.GET("/api/health", updateHandler.HealthCheck)

	router.Run(":8080")
}