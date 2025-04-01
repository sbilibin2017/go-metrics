package repositories_test

import (
	"context"
	"go-metrics/internal/domain"
	"go-metrics/internal/engines"
	"go-metrics/internal/repositories"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type mockDSNGetter2 struct {
	dsn string
}

func (m *mockDSNGetter2) GetDatabaseDSN() string {
	return m.dsn
}

func TestMetricDBFindRepository_Find(t *testing.T) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForSQL("5432/tcp", "pgx", func(host string, port nat.Port) string {
				return "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
			}),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	getter := &mockDSNGetter2{dsn: dsn}

	dbEngine := &engines.DBEngine{}
	err = dbEngine.Open(ctx, getter)
	require.NoError(t, err)
	defer dbEngine.Close()

	_, err = dbEngine.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metrics (
			id TEXT NOT NULL,
			type TEXT NOT NULL,
			delta BIGINT,
			value DOUBLE PRECISION,
			PRIMARY KEY (id, type)
		);
	`)
	require.NoError(t, err)

	repo := repositories.NewMetricDBSaveRepository(dbEngine)

	metrics := []*domain.Metric{
		{ID: "cpu", Type: "gauge", Value: new(float64)},
		{ID: "requests", Type: "counter", Delta: new(int64)},
	}
	*metrics[0].Value = 1.23
	*metrics[1].Delta = 10

	err = repo.Save(ctx, metrics)
	assert.NoError(t, err)

	findRepo := repositories.NewMetricDBFindRepository(dbEngine)
	filters := []*domain.MetricID{
		{ID: "cpu", Type: "gauge"},
		{ID: "requests", Type: "counter"},
	}

	result, err := findRepo.Find(ctx, filters)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, *metrics[0].Value, *result[domain.MetricID{ID: "cpu", Type: "gauge"}].Value)
	assert.Equal(t, *metrics[1].Delta, *result[domain.MetricID{ID: "requests", Type: "counter"}].Delta)
}
