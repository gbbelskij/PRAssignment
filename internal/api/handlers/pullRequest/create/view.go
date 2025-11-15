package pullRequestCreate

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/logger"
	customErrors "PRAssignment/internal/repository/custom_errors"
	"PRAssignment/internal/request"
	"PRAssignment/internal/response"
	service "PRAssignment/internal/service/pullRequest"
	"PRAssignment/pkg"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PullRequestAdder interface {
	AddPullRequest(ctx context.Context, pullRequest *domain.PullRequest) ([]string, error)
}

func Handle(log *slog.Logger, pullRequestAdder PullRequestAdder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.PullRequestCreateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", logger.Err(err))
			c.JSON(http.StatusBadRequest, response.MakeError(
				response.ErrCodeBadRequest,
				"Invalid request",
			))
			return
		}

		req.AuthorId = pkg.ParseOrGenerateUUID(req.AuthorId)
		req.PullRequestId = pkg.ParseOrGenerateUUID(req.PullRequestId)

		pullRequest := service.PullRequestFromRequest(req)
		reviewers, err := pullRequestAdder.AddPullRequest(c.Request.Context(), pullRequest)
		if err != nil {
			if errors.Is(err, customErrors.ErrNotFound) {
				log.Error("no such author", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"No such author",
				))
				return
			}

			if errors.Is(err, customErrors.ErrUniqueViolation) {
				log.Error("PR with such id already exists", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodePRExists,
					"PR with such id already exists",
				))
				return
			}

			if errors.Is(err, customErrors.ErrNoCandidate) {
				log.Error("no candidates", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodeNoCandidate,
					"No candidates",
				))
				return
			}

			log.Error("failed to create pull request", logger.Err(err))
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Internal server error",
			))
			return
		}

		c.JSON(http.StatusCreated, response.MakePullRequestCreateResponse(
			req.PullRequestId,
			req.PullRequestName,
			req.AuthorId,
			domain.PullRequestStatusOpen,
			reviewers,
		))
	}
}
