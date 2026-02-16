package usecase_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/knjname/go-todo-api/internal/domain"
	"github.com/knjname/go-todo-api/internal/usecase"
	"github.com/knjname/go-todo-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestUseCase(repo *mocks.TodoRepository) *usecase.TodoUseCase {
	logger := slog.New(slog.DiscardHandler)
	return usecase.NewTodoUseCase(repo, logger)
}

func TestCreateTodo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil)
		uc := newTestUseCase(repo)

		todo, err := uc.CreateTodo(context.Background(), "Test", "Description")
		require.NoError(t, err)
		assert.Equal(t, "Test", todo.Title)
		assert.Equal(t, "Description", todo.Description)
		assert.False(t, todo.Completed)
		repo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		uc := newTestUseCase(repo)

		_, err := uc.CreateTodo(context.Background(), "", "Description")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrValidation)
	})
}

func TestGetTodo(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		id := uuid.New()
		expected := &domain.Todo{ID: id, Title: "Test"}
		repo.On("GetByID", mock.Anything, id).Return(expected, nil)
		uc := newTestUseCase(repo)

		todo, err := uc.GetTodo(context.Background(), id)
		require.NoError(t, err)
		assert.Equal(t, expected, todo)
	})

	t.Run("not found", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		id := uuid.New()
		repo.On("GetByID", mock.Anything, id).Return(nil, domain.ErrNotFound)
		uc := newTestUseCase(repo)

		_, err := uc.GetTodo(context.Background(), id)
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestListTodos(t *testing.T) {
	repo := mocks.NewTodoRepository(t)
	expected := []domain.Todo{{Title: "A"}, {Title: "B"}}
	repo.On("List", mock.Anything).Return(expected, nil)
	uc := newTestUseCase(repo)

	todos, err := uc.ListTodos(context.Background())
	require.NoError(t, err)
	assert.Len(t, todos, 2)
}

func TestUpdateTodo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		id := uuid.New()
		existing := &domain.Todo{ID: id, Title: "Old", Description: "Old desc"}
		repo.On("GetByID", mock.Anything, id).Return(existing, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil)
		uc := newTestUseCase(repo)

		todo, err := uc.UpdateTodo(context.Background(), id, "New", "New desc")
		require.NoError(t, err)
		assert.Equal(t, "New", todo.Title)
		assert.Equal(t, "New desc", todo.Description)
	})

	t.Run("not found", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		id := uuid.New()
		repo.On("GetByID", mock.Anything, id).Return(nil, domain.ErrNotFound)
		uc := newTestUseCase(repo)

		_, err := uc.UpdateTodo(context.Background(), id, "New", "desc")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewTodoRepository(t)
		id := uuid.New()
		repo.On("GetByID", mock.Anything, id).Return(&domain.Todo{ID: id}, nil)
		repo.On("Delete", mock.Anything, id).Return(nil)
		uc := newTestUseCase(repo)

		err := uc.DeleteTodo(context.Background(), id)
		require.NoError(t, err)
	})
}

func TestCompleteTodo(t *testing.T) {
	repo := mocks.NewTodoRepository(t)
	id := uuid.New()
	existing := &domain.Todo{ID: id, Title: "Task", Completed: false}
	repo.On("GetByID", mock.Anything, id).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil)
	uc := newTestUseCase(repo)

	todo, err := uc.CompleteTodo(context.Background(), id)
	require.NoError(t, err)
	assert.True(t, todo.Completed)
}

func TestCompleteAllTodos(t *testing.T) {
	repo := mocks.NewTodoRepository(t)
	repo.On("CompleteAll", mock.Anything).Return(int64(5), nil)
	uc := newTestUseCase(repo)

	count, err := uc.CompleteAllTodos(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
