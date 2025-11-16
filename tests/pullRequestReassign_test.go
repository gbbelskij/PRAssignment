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

func TestShouldReassignPullRequestHappyPath(t *testing.T) {
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

	teamAddRequest, err := os.ReadFile("static/teamAddRequestMoreMembers.json")
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

	userSetIsActiveRequest, err := os.ReadFile("static/usersSetIsActiveStas.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	_, err = http.Post(server.URL+"/users/setIsActive", "application/json", strings.NewReader(string(userSetIsActiveRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	reassignPullRequest, err := os.ReadFile("static/pullRequestReassignRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/reassign", "application/json", strings.NewReader(string(reassignPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	expectedRequest, err := os.ReadFile("static/pullRequestReassignResponse.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	expectedJSON := string(expectedRequest)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	respBody := string(respBodyBytes)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")
	assert.JSONEq(t, expectedJSON, respBody, "Response JSON should match expected JSON")
}

func TestShouldReturn404WhenPRNotFound(t *testing.T) {
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

	reassignPullRequest, err := os.ReadFile("static/pullRequestReassignRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/reassign", "application/json", strings.NewReader(string(reassignPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code should be 404")
}

func TestShouldReturn409WhenPRMerged(t *testing.T) {
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

	teamAddRequest, err := os.ReadFile("static/teamAddRequestMoreMembers.json")
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

	_, err = http.Post(server.URL+"/pullRequest/merge", "application/json", strings.NewReader(string(mergePullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	reassignPullRequest, err := os.ReadFile("static/pullRequestReassignRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/reassign", "application/json", strings.NewReader(string(reassignPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var errResp ErrorResponse
	err = json.Unmarshal(respBodyBytes, &errResp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Status code should be 409")
	assert.Equal(t, "PR_MERGED", errResp.Error.Code, "Error code should be PR_MERGED")
}

func TestShouldReturn409WhenAuthorNotAssigned(t *testing.T) {
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

	teamAddRequest, err := os.ReadFile("static/teamAddRequestMoreMembers.json")
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

	reassignPullRequest, err := os.ReadFile("static/pullRequestReassignRequestNotAssigned.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/reassign", "application/json", strings.NewReader(string(reassignPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var errResp ErrorResponse
	err = json.Unmarshal(respBodyBytes, &errResp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Status code should be 409")
	assert.Equal(t, "NOT_ASSIGNED", errResp.Error.Code, "Error code should be NOT_ASSIGNED")
}

func TestShouldReturn409WhenNoCandidates(t *testing.T) {
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

	reassignPullRequest, err := os.ReadFile("static/pullRequestReassignRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/pullRequest/reassign", "application/json", strings.NewReader(string(reassignPullRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var errResp ErrorResponse
	err = json.Unmarshal(respBodyBytes, &errResp)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Status code should be 409")
	assert.Equal(t, "NO_CANDIDATE", errResp.Error.Code, "Error code should be NO_CANDIDATE")
}
