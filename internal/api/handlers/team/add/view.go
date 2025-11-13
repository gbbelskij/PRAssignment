package teamAdd

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	"PRAssignment/internal/request"
	service "PRAssignment/internal/service/team"
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamSaver interface {
	SaveTeam(ctx context.Context, team *domain.Team) error
	SaveMembers(ctx context.Context, members []domain.TeamMember) error
}

func Handle(ctx context.Context, log *slog.Logger, teamSaver TeamSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.TeamAddRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := teamSaver.SaveTeam(
			ctx,
			service.TeamFromRequest(req.Team),
		)
		if err != nil {
			log.Error("failed to save team", logger.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save team"})
			return
		}

		log.Info("saved team successfully")

		c.JSON(http.StatusOK, gin.H{"message": "saved team successfully"})
	}
}
