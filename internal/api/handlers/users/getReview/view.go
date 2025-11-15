package usersGetReview

import (
	"PRAssignment/internal/domain"
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PullRequestsGetter interface {
	GetPullRequests(ctx context.Context, userId string) ([]domain.PullRequest, error)
}

func Handle(log *slog.Logger, pullRequestsGetter PullRequestsGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := pkg.ParseOrGenerateUUID(c.Query("user_id"))

		pullRequests, err := pullRequestsGetter.GetPullRequests(c.Request.Context(), userId)
		if err != nil {
			log.Error("failed go get pull requests")
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to get pull requests",
			))
		}

		c.JSON(http.StatusOK, response.MakeUsersGetReviewResponse(
			userId,
			pullRequests,
		))
	}
}
