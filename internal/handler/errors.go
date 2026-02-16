package handler

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/knjname/go-todo-api/internal/domain"
)

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("resource not found", err)
	case errors.Is(err, domain.ErrValidation):
		return huma.Error422UnprocessableEntity("validation failed", err)
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
