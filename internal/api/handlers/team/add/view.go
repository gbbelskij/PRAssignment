package teamAdd

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/request"
	service "PRAssignment/internal/service/team"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamAdder interface {
	SaveTeamWithMembers(ctx context.Context, team *domain.Team, members []domain.TeamMember) (string, error)
}

func Handle(ctx context.Context, log *slog.Logger, teamAdder TeamAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.TeamAddRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		teamId, err := teamAdder.SaveTeamWithMembers(
			c.Request.Context(),
			service.TeamFromRequest(req.Team),
			service.TeamMembersFromRequest(req.Team.Members),
		)

		if err != nil {
			log.Error("failed to save team with members", logger.Err(err))

			if errors.Is(err, customErrors.ErrUniqueViolation) {
				c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("%s already exists", req.Team.TeamName)})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save team"})
			return
		}

		log.Info("saved team with members successfully", "team_id", teamId)
		c.JSON(http.StatusCreated, gin.H{"message": "saved team successfully", "team_id": teamId})
	}
}
