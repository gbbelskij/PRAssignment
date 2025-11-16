package pullRequestReassign

import (
	"PRAssignment/internal/logger"
	customerrors "PRAssignment/internal/repository/customErrors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PullRequestReassignService interface {
	ReassignReviewer(ctx context.Context, pullRequestID string, oldUserID string) (*response.PullRequestReassignResponse, error)
}

func Handle(log *slog.Logger, svc PullRequestReassignService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.PullRequestReassignRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, response.MakeError(
				response.ErrCodeBadRequest,
				"Invalid request",
			))
			return
		}

		req.PullRequestID = pkg.ParseOrGenerateUUID(req.PullRequestID)
		req.OldUserID = pkg.ParseOrGenerateUUID(req.OldUserID)

		resp, err := svc.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				log.Error("no such pr or user", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"No such pr or user",
				))
				return
			}

			if errors.Is(err, customerrors.ErrPrMerged) {
				log.Error("pr merged", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodePRMerged,
					"Pr merged",
				))
				return
			}

			if errors.Is(err, customerrors.ErrNotAssigned) {
				log.Error("not assigned", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodeNotAssigned,
					"Not assigned",
				))
				return
			}

			if errors.Is(err, customerrors.ErrNoCandidate) {
				log.Error("no candidate", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodeNoCandidate,
					"No candidate",
				))
				return
			}

			log.Error("failed to reassign reviewer", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to reassign reviewer",
			))
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
