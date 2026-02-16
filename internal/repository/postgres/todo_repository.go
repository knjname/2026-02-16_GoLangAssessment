package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knjname/go-todo-api/internal/domain"
)

type TodoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoRepository(pool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{pool: pool}
}

func (r *TodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	_, err := r.pool.Exec(ctx, queryInsertTodo,
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.CreatedAt, todo.UpdatedAt,
	)
	return err
}

func (r *TodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	var t domain.Todo
	err := r.pool.QueryRow(ctx, queryGetTodoByID, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepository) List(ctx context.Context) ([]domain.Todo, error) {
	rows, err := r.pool.Query(ctx, queryListTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var t domain.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	tag, err := r.pool.Exec(ctx, queryUpdateTodo,
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, queryDeleteTodo, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TodoRepository) CompleteAll(ctx context.Context) (int64, error) {
	tag, err := r.pool.Exec(ctx, queryCompleteAll)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
