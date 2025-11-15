package pullRequestReassign

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

type ReviewerUpdater interface {
	UpdateReviewer(ctx context.Context, pullRequestId string, oldUserId string) (*domain.PullRequest, string, error)
	GetPullRequestReviewers(ctx context.Context, pullRequestId string) ([]string, error)
}

func Handle(log *slog.Logger, reviewerUpdater ReviewerUpdater) gin.HandlerFunc {
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

		req.PullRequestId = pkg.ParseOrGenerateUUID(req.PullRequestId)
		req.OldUserId = pkg.ParseOrGenerateUUID(req.OldUserId)

		pullRequest, replacedBy, err := reviewerUpdater.UpdateReviewer(c.Request.Context(), req.PullRequestId, req.OldUserId)
		if err != nil {
			if errors.Is(err, customErrors.ErrNotFound) {
				log.Error("no such pr or user", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"No such pr or user",
				))
				return
			}

			if errors.Is(err, customErrors.ErrPrMerged) {
				log.Error("pr merged", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodePRMerged,
					"Pr merged",
				))
				return
			}

			if errors.Is(err, customErrors.ErrNotAssigned) {
				log.Error("not assigned", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodeNotAssigned,
					"Not assigned",
				))
				return
			}

			if errors.Is(err, customErrors.ErrNoCandidate) {
				log.Error("no candidate", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodeNoCandidate,
					"No candidate",
				))
				return
			}
		}

		reviewers, err := reviewerUpdater.GetPullRequestReviewers(c.Request.Context(), req.PullRequestId)
		if err != nil {
			log.Error("failed to get reviewers", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to get reviewers",
			))
			return
		}

		c.JSON(http.StatusOK, response.MakePullRequestReassignResponse(
			pullRequest.PullRequestID,
			pullRequest.PullRequestName,
			pullRequest.AuthorID,
			pullRequest.Status,
			reviewers,
			replacedBy,
		))
	}
}
