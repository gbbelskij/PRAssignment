package tests

import (
	"PRAssignment/internal/app"
	"PRAssignment/internal/container"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsersSetIsActiveHappyPath(t *testing.T) {
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

	setIsActiveRequest, err := os.ReadFile("static/usersSetIsActiveRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/users/setIsActive", "application/json", strings.NewReader(string(setIsActiveRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	expectedData, err := os.ReadFile("static/usersSetIsActiveResponse.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	expectedJSON := string(expectedData)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	respBody := string(respBodyBytes)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")
	assert.JSONEq(t, expectedJSON, respBody, "Response JSON should match expected JSON")
}

func TestShouldReturn404WhenUserNotFound(t *testing.T) {
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

	setIsActiveRequest, err := os.ReadFile("static/usersSetIsActiveRequest.json")
	if err != nil {
		t.Fatalf("failed to read json file: %v", err)
	}

	resp, err := http.Post(server.URL+"/users/setIsActive", "application/json", strings.NewReader(string(setIsActiveRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code should be 404")
}
