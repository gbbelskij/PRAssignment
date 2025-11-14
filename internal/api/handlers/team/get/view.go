package teamGet

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/response"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamGetter interface {
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	GetMembers(ctx context.Context, teamId string) ([]domain.TeamMember, error)
}

func Handle(log *slog.Logger, teamGetter TeamGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		teamName := c.Query("team_name")

		team, err := teamGetter.GetTeam(c.Request.Context(), teamName)
		if err != nil {
			if errors.Is(err, customErrors.ErrNotFound) {
				log.Error("failed to find team", logger.Err(err))
				c.JSON(http.StatusNotFound, response.ErrorResponse{
					Error: response.Error{
						Code:    response.ErrCodeNotFound,
						Message: "team not found",
					},
				})
				return
			}
		}

		members, err := teamGetter.GetMembers(c.Request.Context(), team.TeamID)
		if err != nil {
			log.Error("failed to find members", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: response.Error{
					Code:    response.ErrBadRequest,
					Message: "members not found",
				},
			})
			return
		}

		var responseMembers []response.TeamMember
		for _, m := range members {
			responseMembers = append(responseMembers, response.TeamMember{
				UserId:   m.UserID,
				Username: m.Username,
				IsActive: m.IsActive,
			})
		}

		resp := response.TeamGetResponse{
			TeamName: team.TeamName,
			Members:  responseMembers,
		}

		c.JSON(http.StatusOK, resp)
	}
}
