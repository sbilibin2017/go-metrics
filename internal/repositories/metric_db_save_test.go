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

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_USER": "testuser", "POSTGRES_PASSWORD": "testpassword", "POSTGRES_DB": "testdb"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	require.NoError(t, err)
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	_, err = db.Exec("CREATE TABLE metrics (id TEXT NOT NULL, type TEXT NOT NULL, delta BIGINT, value DOUBLE PRECISION, PRIMARY KEY (id, type));")
	require.NoError(t, err)
	return db, func() { db.Close(); container.Terminate(ctx) }
}

func TestDBSave(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repositories.NewMetricDBSaveRepository(db)
	delta := int64(10)
	value := 20.5
	metrics := []*domain.Metric{{ID: "1", Type: "counter", Delta: &delta}, {ID: "2", Type: "gauge", Value: &value}}
	err := repo.Save(context.Background(), metrics)
	require.NoError(t, err)
	var id, mType string
	var deltaRes sql.NullInt64
	var valueRes sql.NullFloat64
	err = db.QueryRow("SELECT id, type, delta, value FROM metrics WHERE id = '1'").Scan(&id, &mType, &deltaRes, &valueRes)
	require.NoError(t, err)
	assert.Equal(t, "1", id)
	assert.Equal(t, "counter", mType)
	assert.True(t, deltaRes.Valid)
	assert.Equal(t, delta, deltaRes.Int64)
	assert.False(t, valueRes.Valid)
	err = db.QueryRow("SELECT id, type, delta, value FROM metrics WHERE id = '2'").Scan(&id, &mType, &deltaRes, &valueRes)
	require.NoError(t, err)
	assert.Equal(t, "2", id)
	assert.Equal(t, "gauge", mType)
	assert.False(t, deltaRes.Valid)
	assert.True(t, valueRes.Valid)
	assert.Equal(t, value, valueRes.Float64)
}
