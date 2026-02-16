package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/knjname/go-todo-api/internal/usecase"
)

type TodoHandler struct {
	uc *usecase.TodoUseCase
}

func NewTodoHandler(uc *usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

// --- Input/Output types ---

type TodoBody struct {
	ID          uuid.UUID `json:"id" doc:"Todo ID"`
	Title       string    `json:"title" doc:"Todoタイトル"`
	Description string    `json:"description" doc:"詳細説明"`
	Completed   bool      `json:"completed" doc:"完了フラグ"`
	CreatedAt   time.Time `json:"createdAt" doc:"作成日時"`
	UpdatedAt   time.Time `json:"updatedAt" doc:"更新日時"`
}

type CreateTodoInput struct {
	Body struct {
		Title       string `json:"title" maxLength:"200" minLength:"1" doc:"Todoタイトル"`
		Description string `json:"description" doc:"詳細説明"`
	}
}

type CreateTodoOutput struct {
	Body TodoBody
}

type GetTodoInput struct {
	ID uuid.UUID `path:"id" doc:"Todo ID"`
}

type GetTodoOutput struct {
	Body TodoBody
}

type ListTodosOutput struct {
	Body []TodoBody
}

type UpdateTodoInput struct {
	ID   uuid.UUID `path:"id" doc:"Todo ID"`
	Body struct {
		Title       string `json:"title" maxLength:"200" minLength:"1" doc:"Todoタイトル"`
		Description string `json:"description" doc:"詳細説明"`
	}
}

type UpdateTodoOutput struct {
	Body TodoBody
}

type DeleteTodoInput struct {
	ID uuid.UUID `path:"id" doc:"Todo ID"`
}

type CompleteTodoInput struct {
	ID uuid.UUID `path:"id" doc:"Todo ID"`
}

type CompleteTodoOutput struct {
	Body TodoBody
}

type CompleteAllOutput struct {
	Body struct {
		Count int64 `json:"count" doc:"完了にしたTodo数"`
	}
}

// Register registers all todo routes on the huma API.
func (h *TodoHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-todo",
		Method:      http.MethodPost,
		Path:        "/todos",
		Summary:     "Create a new todo",
		Tags:        []string{"Todos"},
	}, h.createTodo)

	huma.Register(api, huma.Operation{
		OperationID: "get-todo",
		Method:      http.MethodGet,
		Path:        "/todos/{id}",
		Summary:     "Get a todo by ID",
		Tags:        []string{"Todos"},
	}, h.getTodo)

	huma.Register(api, huma.Operation{
		OperationID: "list-todos",
		Method:      http.MethodGet,
		Path:        "/todos",
		Summary:     "List all todos",
		Tags:        []string{"Todos"},
	}, h.listTodos)

	huma.Register(api, huma.Operation{
		OperationID: "update-todo",
		Method:      http.MethodPut,
		Path:        "/todos/{id}",
		Summary:     "Update a todo",
		Tags:        []string{"Todos"},
	}, h.updateTodo)

	huma.Register(api, huma.Operation{
		OperationID: "delete-todo",
		Method:      http.MethodDelete,
		Path:        "/todos/{id}",
		Summary:     "Delete a todo",
		Tags:        []string{"Todos"},
	}, h.deleteTodo)

	huma.Register(api, huma.Operation{
		OperationID: "complete-todo",
		Method:      http.MethodPost,
		Path:        "/todos/{id}/complete",
		Summary:     "Mark a todo as complete",
		Tags:        []string{"Todos"},
	}, h.completeTodo)

	huma.Register(api, huma.Operation{
		OperationID: "complete-all-todos",
		Method:      http.MethodPost,
		Path:        "/todos/complete-all",
		Summary:     "Mark all todos as complete",
		Tags:        []string{"Todos"},
	}, h.completeAllTodos)
}

func (h *TodoHandler) createTodo(ctx context.Context, input *CreateTodoInput) (*CreateTodoOutput, error) {
	todo, err := h.uc.CreateTodo(ctx, input.Body.Title, input.Body.Description)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &CreateTodoOutput{Body: TodoBody{
		ID: todo.ID, Title: todo.Title, Description: todo.Description,
		Completed: todo.Completed, CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt,
	}}, nil
}

func (h *TodoHandler) getTodo(ctx context.Context, input *GetTodoInput) (*GetTodoOutput, error) {
	todo, err := h.uc.GetTodo(ctx, input.ID)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &GetTodoOutput{Body: TodoBody{
		ID: todo.ID, Title: todo.Title, Description: todo.Description,
		Completed: todo.Completed, CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt,
	}}, nil
}

func (h *TodoHandler) listTodos(ctx context.Context, _ *struct{}) (*ListTodosOutput, error) {
	todos, err := h.uc.ListTodos(ctx)
	if err != nil {
		return nil, mapDomainError(err)
	}
	body := make([]TodoBody, len(todos))
	for i, t := range todos {
		body[i] = TodoBody{
			ID: t.ID, Title: t.Title, Description: t.Description,
			Completed: t.Completed, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
		}
	}
	return &ListTodosOutput{Body: body}, nil
}

func (h *TodoHandler) updateTodo(ctx context.Context, input *UpdateTodoInput) (*UpdateTodoOutput, error) {
	todo, err := h.uc.UpdateTodo(ctx, input.ID, input.Body.Title, input.Body.Description)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &UpdateTodoOutput{Body: TodoBody{
		ID: todo.ID, Title: todo.Title, Description: todo.Description,
		Completed: todo.Completed, CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt,
	}}, nil
}

func (h *TodoHandler) deleteTodo(ctx context.Context, input *DeleteTodoInput) (*struct{}, error) {
	if err := h.uc.DeleteTodo(ctx, input.ID); err != nil {
		return nil, mapDomainError(err)
	}
	return nil, nil
}

func (h *TodoHandler) completeTodo(ctx context.Context, input *CompleteTodoInput) (*CompleteTodoOutput, error) {
	todo, err := h.uc.CompleteTodo(ctx, input.ID)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &CompleteTodoOutput{Body: TodoBody{
		ID: todo.ID, Title: todo.Title, Description: todo.Description,
		Completed: todo.Completed, CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt,
	}}, nil
}

func (h *TodoHandler) completeAllTodos(ctx context.Context, _ *struct{}) (*CompleteAllOutput, error) {
	count, err := h.uc.CompleteAllTodos(ctx)
	if err != nil {
		return nil, mapDomainError(err)
	}
	out := &CompleteAllOutput{}
	out.Body.Count = count
	return out, nil
}
