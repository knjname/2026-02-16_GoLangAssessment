package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/knjname/go-todo-api/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2 --name=TodoRepository --output=./mocks --outpkg=mocks
type TodoRepository interface {
	Create(ctx context.Context, todo *domain.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	List(ctx context.Context) ([]domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	CompleteAll(ctx context.Context) (int64, error)
}
