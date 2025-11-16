package usersIsActive

import (
	"PRAssignment/internal/logger"
	"PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	SetUserActiveStatus(ctx context.Context, userId string, isActive bool) (*response.UserSetIsActiveResponse, error)
}

func Handle(log *slog.Logger, userService UserService) gin.HandlerFunc {
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

		resp, err := userService.SetUserActiveStatus(
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

			log.Error("failed to update user status", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to update user",
			))
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
