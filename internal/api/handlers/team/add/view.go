package teamAdd

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/request"
	service "PRAssignment/internal/service/team"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamAdder interface {
	SaveTeam(ctx context.Context, team *domain.Team) (string, error)
	SaveTeamMembersBatch(ctx context.Context, members []domain.TeamMember) error
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func Handle(ctx context.Context, log *slog.Logger, teamAdder TeamAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.TeamAddRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := teamAdder.WithTx(c.Request.Context(), func(ctx context.Context) error {
			teamId, err := teamAdder.SaveTeam(
				ctx,
				service.TeamFromRequest(req.Team),
			)
			if err != nil {
				return err
			}
			log.Info("saved team successfully")

			err = teamAdder.SaveTeamMembersBatch(
				ctx,
				service.TeamMembersFromRequest(teamId, req.Team.Members),
			)
			if err != nil {
				return err
			}
			log.Info("saved team members successfully")

			return nil
		})

		if err != nil {
			log.Error("failed to save team or team members", logger.Err(err))
			if errors.Is(err, customErrors.ErrUniqueViolation) {
				c.JSON(http.StatusConflict, gin.H{"message": "team or team member already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save team or team members"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "saved team successfully"})
	}
}
