package pullrequestcreate

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

type pullRequestCreateService interface {
	CreatePullRequest(ctx context.Context, req *request.PullRequestCreateRequest) (*response.PullRequestCreateResponse, error)
}

func Handle(log *slog.Logger, svc pullRequestCreateService) gin.HandlerFunc {
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

		req.AuthorID = pkg.ParseOrGenerateUUID(req.AuthorID)
		req.PullRequestID = pkg.ParseOrGenerateUUID(req.PullRequestID)

		resp, err := svc.CreatePullRequest(c.Request.Context(), &req)
		if err != nil {
			if errors.Is(err, customerrors.ErrNotFound) {
				log.Error("no such author", logger.Err(err))
				c.JSON(http.StatusNotFound, response.MakeError(
					response.ErrCodeNotFound,
					"No such author",
				))
				return
			}

			if errors.Is(err, customerrors.ErrUniqueViolation) {
				log.Error("PR with such id already exists", logger.Err(err))
				c.JSON(http.StatusConflict, response.MakeError(
					response.ErrCodePRExists,
					"PR with such id already exists",
				))
				return
			}

			if errors.Is(err, customerrors.ErrNoCandidate) {
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

		c.JSON(http.StatusCreated, resp)
	}
}
