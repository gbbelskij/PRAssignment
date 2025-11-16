package teamadd

import (
	"PRAssignment/internal/logger"
	customerrors "PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type teamAddService interface {
	AddTeam(ctx context.Context, req *request.TeamAddRequest) (string, error)
}

func Handle(log *slog.Logger, teamAddSvc teamAddService) gin.HandlerFunc {
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
			req.Members[idx].UserID = pkg.ParseOrGenerateUUID(req.Members[idx].UserID)
		}

		teamID, err := teamAddSvc.AddTeam(c.Request.Context(), &req)
		if err != nil {
			if errors.Is(err, customerrors.ErrUniqueViolation) {
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

		log.Info("saved team with members successfully", "team_id", teamID)
		c.JSON(http.StatusCreated, req)
	}
}
