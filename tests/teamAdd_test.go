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

func TestTeamAddHappyPath(t *testing.T) {
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

	resp, err := http.Post(server.URL+"/team/add", "application/json", strings.NewReader(string(teamAddRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	expectedJSON := string(teamAddRequest)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	respBody := string(respBodyBytes)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status code should be 201")
	assert.JSONEq(t, expectedJSON, respBody, "Response JSON should match expected JSON")
}

func TestShoulReturn400WhenTeamExists(t *testing.T) {
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
	resp, err := http.Post(server.URL+"/team/add", "application/json", strings.NewReader(string(teamAddRequest)))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code should be 400")
}
