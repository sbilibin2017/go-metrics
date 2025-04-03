package repositories

import (
	"context"
	"database/sql"
	"go-metrics/internal/domain"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestBuildMetricFindQuery_EmptyFilters(t *testing.T) {
	filters := []*domain.MetricID{}
	query, args := buildMetricFindQuery(filters)
	expectedQuery := "SELECT id, type, delta, value FROM metrics"
	expectedArgs := []any{}
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildMetricFindQuery_OneFilter(t *testing.T) {
	filters := []*domain.MetricID{
		{ID: "metric-1", Type: domain.Counter},
	}
	query, args := buildMetricFindQuery(filters)
	expectedQuery := "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2)"
	expectedArgs := []any{"metric-1", domain.Counter}
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedArgs, args)
}

func TestBuildMetricFindQuery_MultipleFilters(t *testing.T) {
	filters := []*domain.MetricID{
		{ID: "metric-1", Type: domain.Counter},
		{ID: "metric-2", Type: domain.Gauge},
	}
	query, args := buildMetricFindQuery(filters)
	expectedQuery := "SELECT id, type, delta, value FROM metrics WHERE (id = $1 AND type = $2) OR (id = $3 AND type = $4)"
	expectedArgs := []any{"metric-1", domain.Counter, "metric-2", domain.Gauge}
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedArgs, args)
}

func runPostgresContainer2(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
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

func TestFind_Metrics(t *testing.T) {
	ctx := context.Background()

	// Set up the PostgreSQL container and DB connection
	postgresContainer, db, err := runPostgresContainer2(ctx)
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Create repository instance
	repository := NewMetricDBFindRepository(db)

	// Insert test data into the database
	_, err = db.ExecContext(ctx, `
		INSERT INTO metrics (id, type, delta, value) VALUES
		('metric-1', 'counter', 10, 100.5),
		('metric-2', 'gauge', NULL, 200.5),
		('metric-3', 'counter', 20, NULL);
	`)
	require.NoError(t, err)

	// Define filters
	filters := []*domain.MetricID{
		{ID: "metric-1", Type: domain.Counter},
		{ID: "metric-2", Type: domain.Gauge},
	}

	// Call the Find method
	result, err := repository.Find(ctx, filters)
	require.NoError(t, err)

	// Validate the results
	assert.Len(t, result, 2)
	assert.Contains(t, result, domain.MetricID{ID: "metric-1", Type: domain.Counter})
	assert.Contains(t, result, domain.MetricID{ID: "metric-2", Type: domain.Gauge})

	// Check the values for metric-1
	metric1 := result[domain.MetricID{ID: "metric-1", Type: domain.Counter}]
	assert.Equal(t, int64(10), *metric1.Delta)
	assert.Equal(t, 100.5, *metric1.Value)

	// Check the values for metric-2
	metric2 := result[domain.MetricID{ID: "metric-2", Type: domain.Gauge}]
	assert.Nil(t, metric2.Delta)
	assert.Equal(t, 200.5, *metric2.Value)
}

func TestFind_EmptyResults(t *testing.T) {
	ctx := context.Background()

	// Set up the PostgreSQL container and DB connection
	postgresContainer, db, err := runPostgresContainer2(ctx)
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Create repository instance
	repository := NewMetricDBFindRepository(db)

	// Insert test data into the database
	_, err = db.ExecContext(ctx, `
		INSERT INTO metrics (id, type, delta, value) VALUES
		('metric-1', 'counter', 10, 100.5),
		('metric-2', 'gauge', NULL, 200.5);
	`)
	require.NoError(t, err)

	// Define a filter that will return no results
	filters := []*domain.MetricID{
		{ID: "metric-3", Type: domain.Counter},
	}

	// Call the Find method
	result, err := repository.Find(ctx, filters)
	require.NoError(t, err)

	// Validate the results
	assert.Len(t, result, 0)
}
