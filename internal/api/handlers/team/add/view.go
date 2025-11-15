package teamAdd

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	service "PRAssignment/internal/service/team"
	"PRAssignment/pkg"
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

func Handle(log *slog.Logger, teamAdder TeamAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.TeamAddRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, response.MakeError(
				response.ErrCodeBadRequest,
				"Invalid request",
			))
			return
		}

		for idx := range req.Members {
			req.Members[idx].UserId = pkg.ParseOrGenerateUUID(req.Members[idx].UserId)
		}

		teamId, err := teamAdder.SaveTeamWithMembers(
			c.Request.Context(),
			service.TeamFromRequest(req),
			service.TeamMembersFromRequest(req.Members),
		)

		if err != nil {
			log.Error("failed to save team with members", logger.Err(err))

			if errors.Is(err, customErrors.ErrUniqueViolation) {
				log.Error("team already exists", logger.Err(err))
				c.JSON(http.StatusBadRequest, response.MakeError(
					response.ErrCodeTeamExists,
					fmt.Sprintf("%s already exists", req.TeamName),
				))
				return
			}

			log.Error("failed to add team", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to add team",
			))
			return
		}

		log.Info("saved team with members successfully", "team_id", teamId)
		c.JSON(http.StatusCreated, req)
	}
}
