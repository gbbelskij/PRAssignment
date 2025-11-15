package pullRequestMerge

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

type PullRequestUpdater interface {
	UpdatePullRequest(ctx context.Context, pullRequestId string) (*domain.PullRequest, error)
	GetPullRequestReviewers(ctx context.Context, pullRequestId string) ([]string, error)
}

func Handle(log *slog.Logger, pullRequestUpdater PullRequestUpdater) gin.HandlerFunc {
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

		req.PullRequestId = pkg.ParseOrGenerateUUID(req.PullRequestId)

		pullRequest, err := pullRequestUpdater.UpdatePullRequest(c.Request.Context(), req.PullRequestId)
		if err != nil {
			if errors.Is(err, customErrors.ErrNotFound) {
				log.Error("pull request not found", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"Pull request not found",
				))
				return
			}

			log.Error("failed to update pull request", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to update pull request",
			))
			return
		}

		reviewers, err := pullRequestUpdater.GetPullRequestReviewers(c.Request.Context(), req.PullRequestId)
		if err != nil {
			log.Error("failed to find reviewers", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to find reviewers",
			))
		}

		c.JSON(http.StatusOK, response.MakePullRequestMergeResponse(
			pullRequest.PullRequestID,
			pullRequest.PullRequestName,
			pullRequest.AuthorID,
			pullRequest.Status,
			reviewers,
			pullRequest.MergedAt,
		))
	}
}
