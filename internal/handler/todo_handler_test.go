package handler_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/knjname/go-todo-api/internal/domain"
	"github.com/knjname/go-todo-api/internal/handler"
	"github.com/knjname/go-todo-api/internal/usecase"
	"github.com/knjname/go-todo-api/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAPI(t *testing.T) (humatest.TestAPI, *mocks.TodoRepository) {
	t.Helper()
	repo := mocks.NewTodoRepository(t)
	logger := slog.New(slog.DiscardHandler)
	uc := usecase.NewTodoUseCase(repo, logger)
	_, api := humatest.New(t)
	h := handler.NewTodoHandler(uc)
	h.Register(api)
	return api, repo
}

func TestCreateTodo_Handler(t *testing.T) {
	api, repo := setupAPI(t)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil)

	resp := api.Post("/todos", map[string]string{
		"title":       "Test Todo",
		"description": "A test todo",
	})

	assert.Equal(t, http.StatusOK, resp.Code)

	var body handler.TodoBody
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "Test Todo", body.Title)
	assert.Equal(t, "A test todo", body.Description)
	assert.False(t, body.Completed)
}

func TestCreateTodo_Handler_ValidationError(t *testing.T) {
	api, _ := setupAPI(t)

	resp := api.Post("/todos", map[string]string{
		"title":       "",
		"description": "",
	})

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestGetTodo_Handler(t *testing.T) {
	api, repo := setupAPI(t)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&domain.Todo{
		ID:    id,
		Title: "Found",
	}, nil)

	resp := api.Get("/todos/" + id.String())
	assert.Equal(t, http.StatusOK, resp.Code)

	var body handler.TodoBody
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "Found", body.Title)
}

func TestGetTodo_Handler_NotFound(t *testing.T) {
	api, repo := setupAPI(t)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(nil, domain.ErrNotFound)

	resp := api.Get("/todos/" + id.String())
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestListTodos_Handler(t *testing.T) {
	api, repo := setupAPI(t)
	repo.On("List", mock.Anything).Return([]domain.Todo{
		{Title: "A"},
		{Title: "B"},
	}, nil)

	resp := api.Get("/todos")
	assert.Equal(t, http.StatusOK, resp.Code)

	var body []handler.TodoBody
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Len(t, body, 2)
}

func TestDeleteTodo_Handler(t *testing.T) {
	api, repo := setupAPI(t)
	id := uuid.New()
	repo.On("GetByID", mock.Anything, id).Return(&domain.Todo{ID: id}, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)

	resp := api.Delete("/todos/" + id.String())
	assert.Equal(t, http.StatusNoContent, resp.Code)
}
