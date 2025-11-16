package stats

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsService interface {
	GetStats(ctx context.Context) (*domain.Stats, error)
}

func Handle(log *slog.Logger, svc StatsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := svc.GetStats(c.Request.Context())
		if err != nil {
			log.Error("failed to get stats", logger.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get stats",
			})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}
