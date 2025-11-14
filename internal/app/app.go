package app

import (
	teamAdd "PRAssignment/internal/api/handlers/team/add"
	"PRAssignment/internal/container"
	"context"

	"github.com/gin-gonic/gin"
)

type App struct {
	container *container.Container
	router    *gin.Engine
	address   string
}

func NewApp(ctx context.Context, c *container.Container) *App {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	setUpRoutes(ctx, c, router)

	return &App{
		container: c,
		router:    router,
		address:   c.Config.Address,
	}
}

func (a *App) Run(ctx context.Context) error {
	return a.router.Run(a.address)
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
