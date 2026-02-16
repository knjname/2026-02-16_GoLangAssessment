package di

//go:generate go tool kessoku $GOFILE

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knjname/go-todo-api/internal/config"
	"github.com/knjname/go-todo-api/internal/repository/postgres"
	"github.com/knjname/go-todo-api/internal/usecase"
	"github.com/mazrean/kessoku"
)

type APIComponents struct {
	Config  *config.Config
	UseCase *usecase.TodoUseCase
	Logger  *slog.Logger
	Pool    *pgxpool.Pool
}

func NewAPIComponents(cfg *config.Config, uc *usecase.TodoUseCase, logger *slog.Logger, pool *pgxpool.Pool) *APIComponents {
	return &APIComponents{
		Config:  cfg,
		UseCase: uc,
		Logger:  logger,
		Pool:    pool,
	}
}

var _ = kessoku.Inject[*APIComponents]("InitializeAPI",
	kessoku.Provide(config.Load),
	kessoku.Provide(NewLogger),
	kessoku.Async(kessoku.Provide(NewPool)),
	kessoku.Bind[usecase.TodoRepository](kessoku.Provide(postgres.NewTodoRepository)),
	kessoku.Provide(usecase.NewTodoUseCase),
	kessoku.Provide(NewAPIComponents),
)
