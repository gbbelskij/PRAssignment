package tests

import (
	"PRAssignment/internal/app"
	"PRAssignment/internal/container"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestMergeHappyPath(t *testing.T) {
	stg, teardown := setupTestContainer(t)
	defer teardown()

	ctx := context.Background()
	cont := &container.Container{
		Storage: stg,
		Config:  container.NewContainer(ctx).Config,
		Logger:  container.NewContainer(ctx).Logger,
	}
	app := app.NewApp(cont)

	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	teamAddRequest, err := os.ReadFile("static/teamAddRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	_, err = http.Post(server.URL+"/team/add", "application/json", strings.NewReader(string(teamAddRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	createPullRequest, err := os.ReadFile("static/pullRequestCreateRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	_, err = http.Post(server.URL+"/pullRequest/create", "application/json", strings.NewReader(string(createPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	mergePullRequest, err := os.ReadFile("static/pullRequestMergeRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/merge", "application/json", strings.NewReader(string(mergePullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	expectedJSON, err := os.ReadFile("static/pullRequestMergeResponse.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	var got map[string]interface{}
	var expected map[string]interface{}

	gotJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	json.Unmarshal(gotJSON, &got)
	json.Unmarshal(expectedJSON, &expected)

	got["pr"].(map[string]interface{})["merged_at"] = "anystring"
	expected["pr"].(map[string]interface{})["merged_at"] = "anystring"

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code shoud be 200")
	assert.Equal(t, expected, got, "Response JSON should match expected JSON")
}

func TestShouldReturn404WhenPrNotFound(t *testing.T) {
	stg, teardown := setupTestContainer(t)
	defer teardown()

	ctx := context.Background()
	cont := &container.Container{
		Storage: stg,
		Config:  container.NewContainer(ctx).Config,
		Logger:  container.NewContainer(ctx).Logger,
	}
	app := app.NewApp(cont)

	server := httptest.NewServer(app.GetRouter())
	defer server.Close()

	mergePullRequest, err := os.ReadFile("static/pullRequestMergeRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/merge", "application/json", strings.NewReader(string(mergePullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code shoud be 404")
}
