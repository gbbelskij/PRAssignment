package teamGet

import (
	"PRAssignment/internal/logger"
	customerrors "PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/response"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamGetService interface {
	GetTeamWithMembers(ctx context.Context, teamName string) (*response.TeamGetResponse, error)
}

func Handle(log *slog.Logger, svc TeamGetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamName := c.Query("team_name")

		resp, err := svc.GetTeamWithMembers(c.Request.Context(), teamName)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				log.Error("no such team", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"Team not found",
				))
				return
			}
			log.Error("failed to get team or members", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to get team or members",
			))
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
