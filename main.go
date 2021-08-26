package main

import (
	"context"
	"fmt"
	"github.com/XiovV/docker_control/controller"
	"github.com/XiovV/docker_control/handlers"
	"github.com/docker/docker/api/types"
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

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
		fmt.Println(container.Names)
	}

	//containerId := os.Args[1]
	//imageTag := os.Args[2]
	//
	//err = dockerController.UpdateContainer(containerId, imageTag)
	//
	//if err != nil {
	//	fmt.Println("couldn't update container:", err)
	//}

	updateHandler := handlers.NewUpdateHandler(dockerController)

	router := gin.Default()
	router.POST("/api/containers/update", updateHandler.HandleContainerUpdate)

	router.Run(":8080")
}