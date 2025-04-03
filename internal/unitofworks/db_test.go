package unitofworks

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createPostgresContainer(ctx context.Context) (testcontainers.Container, string, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, "", "", err
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", "", err
	}
	return postgresContainer, host, port.Port(), nil
}

func TestDBUnitOfWork(t *testing.T) {
	ctx := context.Background()
	container, host, port, err := createPostgresContainer(ctx)
	require.NoError(t, err)
	defer func() {
		err := container.Terminate(ctx)
		require.NoError(t, err)
	}()
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port)
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	defer db.Close()
	uow := NewDBUnitOfWork(db)
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)`)
	require.NoError(t, err)
	err = uow.Do(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users (name) VALUES ($1)`, "Alice")
		return err
	})
	require.NoError(t, err)
	var name string
	err = db.QueryRowContext(ctx, `SELECT name FROM users WHERE name = $1`, "Alice").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "Alice", name)
	uow.Do(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users (name) VALUES ($1)`, "Bob")
		if err != nil {
			return err
		}
		return fmt.Errorf("force rollback")
	})
	var count int
	err = db.QueryRowContext(ctx, `SELECT count(*) FROM users WHERE name = $1`, "Bob").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
