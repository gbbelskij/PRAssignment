package main

import (
	"PRAssignment/internal/container"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	container := container.NewContainer()

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	setUpRoutes(router)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	container.Logger.Info("received shutdown signal")

	container.Storage.Close()

	container.Logger.Info("shutting down gracefully")
}

func setUpRoutes(router *gin.Engine) {
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add")
		teamGroup.GET("/get")
	}

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive")
		usersGroup.GET("/getReview")
	}

	prGroup := router.Group("/pullRequest")
	{
		prGroup.POST("/create")
		prGroup.POST("/merge")
		prGroup.POST("/reassign")
	}
}
