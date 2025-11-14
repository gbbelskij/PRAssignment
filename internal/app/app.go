package app

import (
	teamAdd "PRAssignment/internal/api/handlers/team/add"
	teamGet "PRAssignment/internal/api/handlers/team/get"
	"PRAssignment/internal/container"
	"context"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.yaml"),
	))

	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", teamAdd.Handle(container.Logger, container.Storage))
		teamGroup.GET("/get", teamGet.Handle(container.Logger, container.Storage))
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
