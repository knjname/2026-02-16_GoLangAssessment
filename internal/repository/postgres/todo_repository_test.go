package postgres_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knjname/go-todo-api/internal/domain"
	"github.com/knjname/go-todo-api/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")

	container, err := pgcontainer.Run(ctx,
		"postgres:16-alpine",
		pgcontainer.WithDatabase("test"),
		pgcontainer.WithUsername("test"),
		pgcontainer.WithPassword("test"),
		pgcontainer.WithInitScripts(filepath.Join(migrationsDir, "001_create_todos_table.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestTodoRepository_Create_and_GetByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)
	ctx := context.Background()

	todo, err := domain.NewTodo("Integration Test", "Testing with real DB")
	require.NoError(t, err)

	err = repo.Create(ctx, todo)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.Equal(t, todo.ID, got.ID)
	assert.Equal(t, todo.Title, got.Title)
	assert.Equal(t, todo.Description, got.Description)
	assert.Equal(t, todo.Completed, got.Completed)
}

func TestTodoRepository_GetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)

	_, err := repo.GetByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTodoRepository_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)
	ctx := context.Background()

	for i := range 3 {
		todo, _ := domain.NewTodo("Todo "+string(rune('A'+i)), "")
		require.NoError(t, repo.Create(ctx, todo))
	}

	todos, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, todos, 3)
}

func TestTodoRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)
	ctx := context.Background()

	todo, _ := domain.NewTodo("Original", "desc")
	require.NoError(t, repo.Create(ctx, todo))

	require.NoError(t, todo.UpdateTitle("Updated"))
	todo.MarkComplete()
	require.NoError(t, repo.Update(ctx, todo))

	got, err := repo.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", got.Title)
	assert.True(t, got.Completed)
}

func TestTodoRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)
	ctx := context.Background()

	todo, _ := domain.NewTodo("To delete", "")
	require.NoError(t, repo.Create(ctx, todo))
	require.NoError(t, repo.Delete(ctx, todo.ID))

	_, err := repo.GetByID(ctx, todo.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTodoRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)

	err := repo.Delete(context.Background(), uuid.New())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTodoRepository_CompleteAll(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewTodoRepository(pool)
	ctx := context.Background()

	for i := range 3 {
		todo, _ := domain.NewTodo("Todo "+string(rune('A'+i)), "")
		require.NoError(t, repo.Create(ctx, todo))
	}

	count, err := repo.CompleteAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)

	todos, _ := repo.List(ctx)
	for _, td := range todos {
		assert.True(t, td.Completed)
	}

	// Running again should affect 0 rows
	count, err = repo.CompleteAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
