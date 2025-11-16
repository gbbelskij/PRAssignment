package tests

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"PRAssignment/internal/repository/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func applyMigrations(t *testing.T, dbConnString string, migrationsDir string) {
	db, err := sql.Open("pgx", dbConnString)
	if err != nil {
		t.Fatalf("failed to open DB connection for migrations: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations dir: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			path := filepath.Join(migrationsDir, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", path, err)
			}
			_, err = db.Exec(string(content))
			if err != nil {
				t.Fatalf("Failed to exec migration %s: %v", file.Name(), err)
			}
			t.Logf("Applied migration: %s", file.Name())
		}
	}
}

func setupTestContainer(t *testing.T) (*storage.Storage, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if pgContainer != nil {
			logs, errLogs := pgContainer.Logs(ctx)
			if errLogs == nil {
				logData, _ := io.ReadAll(logs)
				t.Logf("Container logs:\n%s", string(logData))
			}
		}
		t.Fatalf("failed to start container: %v", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	t.Logf("Connecting to: %s", dsn)

	maxRetries := 10
	var db *sql.DB
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("pgx", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				t.Logf("Connected to database successfully")
				break
			}
			db.Close()
		}
		if i < maxRetries-1 {
			t.Logf("Retry %d/%d: waiting for database...", i+1, maxRetries)
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	db.Close()

	os.Setenv("POSTGRES_CONN_STRING", dsn)
	os.Setenv("CONFIG_PATH", "../configs/config.yaml")

	applyMigrations(t, dsn, "../migrations")

	stg, err := storage.NewStorage(ctx)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	teardown := func() {
		stg.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
		os.Unsetenv("POSTGRES_CONN_STRING")
	}
	return stg, teardown
}
