package engines_test

import (
	"context"
	"go-metrics/internal/engines"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type mockDSNGetter struct {
	dsn string
}

func (m *mockDSNGetter) GetDatabaseDSN() string {
	return m.dsn
}

func TestDBEngine_Open(t *testing.T) {
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
	getter := &mockDSNGetter{dsn: dsn}

	dbEngine := &engines.DBEngine{}
	err = dbEngine.Open(ctx, getter)
	require.NoError(t, err)
	defer dbEngine.Close()

	assert.NotNil(t, dbEngine.DB)
	err = dbEngine.Ping()
	assert.NoError(t, err)
}
