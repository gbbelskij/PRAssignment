package pullRequestMerge

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

type PullRequestMergeService interface {
	MergePullRequest(ctx context.Context, pullRequestID string) (*response.PullRequestMergeResponse, error)
}

func Handle(log *slog.Logger, svc PullRequestMergeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.PullRequestMergeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, response.MakeError(
				response.ErrCodeBadRequest,
				"Invalid request",
			))
			return
		}

		req.PullRequestID = pkg.ParseOrGenerateUUID(req.PullRequestID)

		resp, err := svc.MergePullRequest(c.Request.Context(), req.PullRequestID)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				log.Error("pull request not found", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"Pull request not found",
				))
				return
			}

			log.Error("failed to merge pull request", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to merge pull request",
			))
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
