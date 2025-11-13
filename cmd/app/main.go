package main

import (
	teamAdd "PRAssignment/internal/api/handlers/team/add"
	"PRAssignment/internal/container"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	container := container.NewContainer()

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	setUpRoutes(ctx, container, router)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	container.Logger.Info("received shutdown signal")

	cancel()
	container.Storage.Close()

	container.Logger.Info("shutting down gracefully")
}

func setUpRoutes(ctx context.Context, container *container.Container, router *gin.Engine) {
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", teamAdd.Handle(ctx, container.Logger, container.Storage))
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
