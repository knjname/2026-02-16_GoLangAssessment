package domain

import (
	"time"

	"github.com/google/uuid"
)

const MaxTitleLength = 200

type Todo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewTodo(title, description string) (*Todo, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Todo{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (t *Todo) UpdateTitle(title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}
	t.Title = title
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Todo) UpdateDescription(description string) {
	t.Description = description
	t.UpdatedAt = time.Now().UTC()
}

func (t *Todo) MarkComplete() {
	t.Completed = true
	t.UpdatedAt = time.Now().UTC()
}

func validateTitle(title string) error {
	if title == "" {
		return NewValidationError("title", "must not be empty")
	}
	if len([]rune(title)) > MaxTitleLength {
		return NewValidationError("title", "must not exceed 200 characters")
	}
	return nil
}
