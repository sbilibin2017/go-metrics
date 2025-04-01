package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"go-metrics/internal/domain"
	"go-metrics/internal/repositories"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB2(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	_, err = db.Exec(`
		CREATE TABLE metrics (
			id TEXT NOT NULL,
			type TEXT NOT NULL,
			delta BIGINT,
			value DOUBLE PRECISION,
			PRIMARY KEY (id, type)
		);
	`)
	require.NoError(t, err)
	cleanup := func() {
		db.Close()
		container.Terminate(ctx)
	}
	return db, cleanup
}

func TestFind(t *testing.T) {
	db, cleanup := setupTestDB2(t)
	defer cleanup()
	repo := repositories.NewMetricDBFindRepository(db)
	_, err := db.Exec(`
		INSERT INTO metrics (id, type, delta, value) VALUES
		('1', 'counter', 10, NULL),
		('2', 'gauge', NULL, 20.5),
		('3', 'counter', 5, NULL),
		('4', 'gauge', NULL, 30.5);
	`)
	require.NoError(t, err)
	filters := []*domain.MetricID{
		{ID: "1", Type: "counter"},
		{ID: "2", Type: "gauge"},
	}
	result, err := repo.Find(context.Background(), filters)
	require.NoError(t, err)
	assert.Len(t, result, 2, "Expected to find 2 metrics")
	metric1, ok := result[domain.MetricID{ID: "1", Type: "counter"}]
	assert.True(t, ok)
	assert.Equal(t, int64(10), *metric1.Delta)
	assert.Nil(t, metric1.Value)
	metric2, ok := result[domain.MetricID{ID: "2", Type: "gauge"}]
	assert.True(t, ok)
	assert.Nil(t, metric2.Delta)
	assert.Equal(t, 20.5, *metric2.Value)
}
