package usersGetReview

import (
	"PRAssignment/internal/response"
	"PRAssignment/pkg"
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetReviewService interface {
	GetUserPullRequests(ctx context.Context, userId string) (*response.UsersGetReviewResponse, error)
}

func Handle(log *slog.Logger, getReviewService GetReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := pkg.ParseOrGenerateUUID(c.Query("user_id"))

		resp, err := getReviewService.GetUserPullRequests(c.Request.Context(), userId)
		if err != nil {
			log.Error("failed to get pull requests", "error", err)
			c.JSON(http.StatusInternalServerError, response.MakeError(
				response.ErrCodeInternalServerError,
				"Failed to get pull requests",
			))
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
