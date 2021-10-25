package main

import (
	"fmt"
	"github.com/XiovV/dokkup-agent/app"
	"github.com/XiovV/dokkup-agent/config"
	"github.com/XiovV/dokkup-agent/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dockerController := controller.New()

	cfg, _, err := config.New("config.json")
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
