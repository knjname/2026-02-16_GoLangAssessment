package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/knjname/go-todo-api/internal/domain"
)

type TodoUseCase struct {
	repo   TodoRepository
	logger *slog.Logger
}

func NewTodoUseCase(repo TodoRepository, logger *slog.Logger) *TodoUseCase {
	return &TodoUseCase{repo: repo, logger: logger}
}

func (uc *TodoUseCase) CreateTodo(ctx context.Context, title, description string) (*domain.Todo, error) {
	todo, err := domain.NewTodo(title, description)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}

	uc.logger.InfoContext(ctx, "todo created", slog.String("id", todo.ID.String()))
	return todo, nil
}

func (uc *TodoUseCase) GetTodo(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get todo: %w", err)
	}
	return todo, nil
}

func (uc *TodoUseCase) ListTodos(ctx context.Context) ([]domain.Todo, error) {
	todos, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	return todos, nil
}

func (uc *TodoUseCase) UpdateTodo(ctx context.Context, id uuid.UUID, title, description string) (*domain.Todo, error) {
	todo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get todo for update: %w", err)
	}

	if err := todo.UpdateTitle(title); err != nil {
		return nil, err
	}
	todo.UpdateDescription(description)

	if err := uc.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}

	uc.logger.InfoContext(ctx, "todo updated", slog.String("id", id.String()))
	return todo, nil
}

func (uc *TodoUseCase) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get todo for delete: %w", err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	uc.logger.InfoContext(ctx, "todo deleted", slog.String("id", id.String()))
	return nil
}

func (uc *TodoUseCase) CompleteTodo(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get todo for complete: %w", err)
	}

	todo.MarkComplete()

	if err := uc.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("complete todo: %w", err)
	}

	uc.logger.InfoContext(ctx, "todo completed", slog.String("id", id.String()))
	return todo, nil
}

func (uc *TodoUseCase) CompleteAllTodos(ctx context.Context) (int64, error) {
	count, err := uc.repo.CompleteAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("complete all todos: %w", err)
	}

	uc.logger.InfoContext(ctx, "all todos completed", slog.Int64("count", count))
	return count, nil
}
