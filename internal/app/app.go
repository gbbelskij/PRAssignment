package app

import (
	pullRequestCreate "PRAssignment/internal/api/handlers/pullRequest/create"
	pullRequestMerge "PRAssignment/internal/api/handlers/pullRequest/merge"
	pullRequestReassign "PRAssignment/internal/api/handlers/pullRequest/reassign"
	"PRAssignment/internal/api/handlers/stats"
	teamAdd "PRAssignment/internal/api/handlers/team/add"
	teamGet "PRAssignment/internal/api/handlers/team/get"
	usersGetReview "PRAssignment/internal/api/handlers/users/getReview"
	usersIsActive "PRAssignment/internal/api/handlers/users/setIsActive"
	"PRAssignment/internal/container"
	"PRAssignment/internal/logger"
	pullRequestService "PRAssignment/internal/service/pullRequest"
	statsService "PRAssignment/internal/service/stats"
	teamService "PRAssignment/internal/service/team"
	userService "PRAssignment/internal/service/users"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	container *container.Container
	router    *gin.Engine
	address   string
}

func NewApp(c *container.Container) *App {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	setUpRoutes(c, router)

	return &App{
		container: c,
		router:    router,
		address:   c.Config.Address,
	}
}

func (a *App) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    a.address,
		Handler: a.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.container.Logger.Error("server error", logger.Err(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.container.Logger.Error("server shutdown error", logger.Err(err))
		return err
	}

	a.container.Logger.Info("server shutdown gracefully")
	return nil
}

func (a *App) GetRouter() *gin.Engine {
	return a.router
}

func setUpRoutes(container *container.Container, router *gin.Engine) {

	setIsActiveService := userService.NewSetIsActiveService(container.Storage)
	getReviewService := userService.NewGetReviewService(container.Storage)
	teamAddService := teamService.NewTeamAddService(container.Storage)
	teamGetService := teamService.NewTeamGetService(container.Storage)
	pullRequestAddService := pullRequestService.NewPullRequestCreateService(container.Storage)
	pullRequestMergeService := pullRequestService.NewPullRequestMergeService(container.Storage)
	pullRequestReassignService := pullRequestService.NewPullRequestReassignService(container.Storage)
	statsService := statsService.NewStatsService(container.Storage)

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.yaml"),
	))

	router.GET("/stats", stats.Handle(container.Logger, statsService))
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", teamAdd.Handle(container.Logger, teamAddService))
		teamGroup.GET("/get", teamGet.Handle(container.Logger, teamGetService))
	}

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive", usersIsActive.Handle(container.Logger, setIsActiveService))
		usersGroup.GET("/getReview", usersGetReview.Handle(container.Logger, getReviewService))
	}

	prGroup := router.Group("/pullRequest")
	{
		prGroup.POST("/create", pullRequestCreate.Handle(container.Logger, pullRequestAddService))
		prGroup.POST("/merge", pullRequestMerge.Handle(container.Logger, pullRequestMergeService))
		prGroup.POST("/reassign", pullRequestReassign.Handle(container.Logger, pullRequestReassignService))
	}
}
