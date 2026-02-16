package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/knjname/go-todo-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
		errField    string
	}{
		{
			name:        "valid todo",
			title:       "Buy groceries",
			description: "Milk, eggs, bread",
			wantErr:     false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "some description",
			wantErr:     true,
			errField:    "title",
		},
		{
			name:        "title too long",
			title:       strings.Repeat("a", 201),
			description: "",
			wantErr:     true,
			errField:    "title",
		},
		{
			name:        "title at max length",
			title:       strings.Repeat("a", 200),
			description: "",
			wantErr:     false,
		},
		{
			name:        "empty description is valid",
			title:       "Valid title",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := domain.NewTodo(tt.title, tt.description)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrValidation))

				var ve *domain.ValidationError
				require.True(t, errors.As(err, &ve))
				assert.Equal(t, tt.errField, ve.Field)
				assert.Nil(t, todo)
			} else {
				require.NoError(t, err)
				require.NotNil(t, todo)
				assert.Equal(t, tt.title, todo.Title)
				assert.Equal(t, tt.description, todo.Description)
				assert.False(t, todo.Completed)
				assert.NotEmpty(t, todo.ID)
				assert.False(t, todo.CreatedAt.IsZero())
				assert.False(t, todo.UpdatedAt.IsZero())
			}
		})
	}
}

func TestTodo_UpdateTitle(t *testing.T) {
	todo, err := domain.NewTodo("Original", "desc")
	require.NoError(t, err)

	err = todo.UpdateTitle("Updated")
	require.NoError(t, err)
	assert.Equal(t, "Updated", todo.Title)

	err = todo.UpdateTitle("")
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
}

func TestTodo_MarkComplete(t *testing.T) {
	todo, err := domain.NewTodo("Task", "desc")
	require.NoError(t, err)
	assert.False(t, todo.Completed)

	todo.MarkComplete()
	assert.True(t, todo.Completed)
}

func TestValidationError_Unwrap(t *testing.T) {
	ve := domain.NewValidationError("field", "msg")
	assert.True(t, errors.Is(ve, domain.ErrValidation))
	assert.Contains(t, ve.Error(), "field")
	assert.Contains(t, ve.Error(), "msg")
}
