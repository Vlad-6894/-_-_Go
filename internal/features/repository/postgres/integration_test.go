//go:build integration

package postgres_repository

import (
	"context"
	"fmt"
	"os"
	"security/internal/core/domains"
	core_postgres_pool "security/internal/core/repository/postgres/pool"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresImage    = "postgres:18.1-bookworm"
	testDatabaseName = "security_test"
	testUser         = "test_user"
	testPassword     = "test123"
	log              = "database system is ready to accept connections"
	localHost        = "127.0.0.1"
)

func TestDatabase_integration(t *testing.T) {
	ctx := t.Context()

	pgContainer, err := postgres.Run(
		ctx,
		postgresImage,
		postgres.WithDatabase(testDatabaseName),
		postgres.WithPassword(testPassword),
		postgres.WithUsername(testUser),
		testcontainers.WithWaitStrategy(wait.ForLog(log).WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("fail to create container: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			fmt.Println("fail to close container", err)
		}
	}()

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	os.Setenv("POSTGRES_HOST", localHost)
	os.Setenv("POSTGRES_PORT", mappedPort.Port())
	os.Setenv("POSTGRES_USER", testUser)
	os.Setenv("POSTGRES_PASSWORD", testPassword)
	os.Setenv("POSTGRES_DB", testDatabaseName)
	os.Setenv("POSTGRES_TIMEOUT", "30s")

	defer func() {
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_DB")
		os.Unsetenv("POSTGRES_TIMEOUT")
	}()

	pool, err := core_postgres_pool.NewPostgresConnectionPool(
		context.Background(),
		core_postgres_pool.NewPostgresConnectionConfigMust(),
	)
	if err != nil {
		t.Fatalf("fail to get pool: %v, port: %v", err, mappedPort.Port())
	}
	defer pool.Close()

	sqlReqInit := `
	CREATE SCHEMA security;

	CREATE TABLE security.dangers (
		id UUID PRIMARY KEY,
		version BIGINT NOT NULL DEFAULT 1 CHECK(version>0),
		name VARCHAR(1000) NOT NULL CHECK(char_length(name) BETWEEN 1 AND 1000),
		md TEXT NOT NULL UNIQUE CHECK(md ~ '^[a-fA-F0-9]{32}$'),
		danger_type VARCHAR(1000) NOT NULL CHECK(char_length(danger_type) BETWEEN 1 AND 1000),
		danger_level VARCHAR(1000) NOT NULL CHECK(danger_level  IN ('LOW', 'MEDIUM', 'HIGH')),
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP
	);
	`

	contextWithTime, cancel := context.WithTimeout(ctx, pool.GetTimeot())
	defer cancel()

	if _, err := pool.Exec(contextWithTime, sqlReqInit); err != nil {
		t.Fatalf("fail to create tables: %v", err)
	}

	repository := NewDangersPostgresRepository(pool)

	t.Run("integration test", func(t *testing.T) {
		danger := domains.NewUnitializedDanger("MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH")

		_, err := repository.CreateDanger(t.Context(), danger)
		if err != nil {
			t.Errorf("fail to create danger: %v", err)
		}

		if _, err := repository.GetDanger(t.Context(), "d41d8cd98f00b204e9800998ecf8427e"); err != nil {
			t.Errorf("fail to get danger: %v", err)
		}

		after := time.Now().Add(5 * time.Second)

		putDanger := domains.NewDanger(uuid.New(), 1, "MyDoom", "d41d8cd98f00b204e9800998ecf8427e", "Virus", "HIGH", time.Now(), &after)

		if _, err := repository.UpdateDanger(t.Context(), putDanger, "d41d8cd98f00b204e9800998ecf8427e"); err != nil {
			t.Errorf("fail to get danger: %v", err)
		}
	})
}
