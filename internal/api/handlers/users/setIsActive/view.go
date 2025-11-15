package usersIsActive

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUpdateTeamGetter interface {
	UpdateUser(ctx context.Context, userId string, isActive bool) (*domain.TeamMember, error)
	GetTeamNameById(ctx context.Context, teamId string) (string, error)
}

func Handle(log *slog.Logger, userUpdaterTeamGetter UserUpdateTeamGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.UserSetIsActiveRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, response.MakeError(
				response.ErrCodeBadRequest,
				"Invalid request",
			))
			return
		}

		req.UserId = pkg.ParseOrGenerateUUID(req.UserId)

		teamMember, err := userUpdaterTeamGetter.UpdateUser(
			c.Request.Context(),
			req.UserId,
			req.IsActive,
		)
		if err != nil {
			if errors.Is(err, customErrors.ErrNotFound) {
				log.Error("no such user", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"User not found",
				))
				return
			}

			log.Error("failed to find user", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to find user",
			))
			return
		}

		teamName, err := userUpdaterTeamGetter.GetTeamNameById(c.Request.Context(), teamMember.TeamID)
		if err != nil {
			log.Error("failed to find team name", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to find team",
			))
			return
		}

		c.JSON(http.StatusOK, response.MakeUserSetIsActiveResponse(
			teamMember.UserID,
			teamMember.Username,
			teamName,
			teamMember.IsActive,
		))
	}
}
