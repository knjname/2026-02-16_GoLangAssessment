package di

//go:generate go tool kessoku $GOFILE

import (
	"database/sql"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knjname/go-todo-api/internal/config"
	"github.com/knjname/go-todo-api/internal/repository/postgres"
	"github.com/knjname/go-todo-api/internal/usecase"
	"github.com/mazrean/kessoku"
)

type BatchComponents struct {
	Config  *config.Config
	UseCase *usecase.TodoUseCase
	Logger  *slog.Logger
	Pool    *pgxpool.Pool
	DB      *sql.DB
}

func NewBatchComponents(cfg *config.Config, uc *usecase.TodoUseCase, logger *slog.Logger, pool *pgxpool.Pool, db *sql.DB) *BatchComponents {
	return &BatchComponents{
		Config:  cfg,
		UseCase: uc,
		Logger:  logger,
		Pool:    pool,
		DB:      db,
	}
}

var _ = kessoku.Inject[*BatchComponents]("InitializeBatch",
	kessoku.Provide(config.Load),
	kessoku.Provide(NewLogger),
	kessoku.Async(kessoku.Provide(NewPool)),
	kessoku.Async(kessoku.Provide(NewStdDB)),
	kessoku.Bind[usecase.TodoRepository](kessoku.Provide(postgres.NewTodoRepository)),
	kessoku.Provide(usecase.NewTodoUseCase),
	kessoku.Provide(NewBatchComponents),
)
