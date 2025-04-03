package repositories

import (
	"context"
	"database/sql"
	"go-metrics/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func runPostgresContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	// Set up PostgreSQL container for testing
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(1), // Wait for the log to appear at least once
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	// Get PostgreSQL container connection info
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	// Create DB connection string
	dbURL := "postgres://testuser:testpassword@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	// Retry logic for establishing a connection
	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("pgx", dbURL)
		if err != nil {
			time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
			continue
		}

		// Attempt to ping the database
		err = db.PingContext(ctx)
		if err == nil {
			break
		}

		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}

	if err != nil {
		return nil, nil, err
	}

	// Set up database schema
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS metrics (
		id TEXT NOT NULL,
		type TEXT NOT NULL,
		delta INT,
		value FLOAT,
		PRIMARY KEY (id, type)
	);
	`)
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, db, nil
}

func TestSave_Metrics(t *testing.T) {
	ctx := context.Background()

	// Set up the PostgreSQL container and DB connection
	postgresContainer, db, err := runPostgresContainer(ctx)
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Create repository instance
	repository := NewMetricDBSaveRepository(db)

	// Prepare test data using the new structure
	metrics := []*domain.Metric{
		{
			MetricID: domain.MetricID{
				ID:   "metric-1",
				Type: domain.Counter,
			},
			Delta: nil,
			Value: nil,
		},
		{
			MetricID: domain.MetricID{
				ID:   "metric-2",
				Type: domain.Gauge,
			},
			Delta: nil,
			Value: nil,
		},
	}

	// Save the metrics
	err = repository.Save(ctx, metrics)
	require.NoError(t, err)

}
