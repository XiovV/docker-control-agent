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

	cfg, err := config.New("config_test.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully loaded config")

	app := app.New(dockerController, cfg)

	router := app.Router()

	fmt.Println("agent is listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
