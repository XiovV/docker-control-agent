package main

import (
	"fmt"
	"github.com/XiovV/docker_control/app"
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dockerController := controller.New()

	cfg := config.New()
	fmt.Println("Successfully loaded config")

	app := app.New(dockerController, cfg)

	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use(app.Authenticate())
	{
		v1.GET("/containers/image/:containerName", app.GetContainerImage)

		v1.PUT("/images/pull", app.PullImage)
		v1.PUT("/containers/update", app.UpdateContainer)
		v1.PUT("/containers/rollback", app.RollbackContainer)
	}

	fmt.Println("agent is listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
