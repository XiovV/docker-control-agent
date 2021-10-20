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
	router.Use(app.Authenticate())

	router.PUT("/v1/images/pull", app.PullImage)
	router.PUT("/v1/containers/update", app.UpdateContainer)

	fmt.Println("agent is listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}