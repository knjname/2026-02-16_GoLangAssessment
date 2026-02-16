package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/knjname/go-todo-api/internal/config"
	"github.com/knjname/go-todo-api/internal/repository/postgres"
	"github.com/knjname/go-todo-api/internal/usecase"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "batch",
		Short: "Todo API batch tool",
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
	}

	migrateUpCmd := &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			db, err := sql.Open("pgx", cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("open database: %w", err)
			}
			defer func() { _ = db.Close() }()

			if err := goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("set dialect: %w", err)
			}
			if err := goose.Up(db, "migrations"); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}
			logger.Info("migrations applied successfully")
			return nil
		},
	}

	migrateDownCmd := &cobra.Command{
		Use:   "down",
		Short: "Roll back the last migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			db, err := sql.Open("pgx", cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("open database: %w", err)
			}
			defer func() { _ = db.Close() }()

			if err := goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("set dialect: %w", err)
			}
			if err := goose.Down(db, "migrations"); err != nil {
				return fmt.Errorf("migrate down: %w", err)
			}
			logger.Info("migration rolled back successfully")
			return nil
		},
	}

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("connect to database: %w", err)
			}
			defer pool.Close()

			repo := postgres.NewTodoRepository(pool)
			uc := usecase.NewTodoUseCase(repo, logger)

			todos, err := uc.ListTodos(ctx)
			if err != nil {
				return fmt.Errorf("list todos: %w", err)
			}

			if len(todos) == 0 {
				fmt.Println("No todos found.")
				return nil
			}

			for _, t := range todos {
				status := "[ ]"
				if t.Completed {
					status = "[x]"
				}
				fmt.Printf("%s %s %s\n", status, t.ID, t.Title)
			}
			return nil
		},
	}

	completeAllCmd := &cobra.Command{
		Use:   "complete-all",
		Short: "Mark all todos as complete",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("connect to database: %w", err)
			}
			defer pool.Close()

			repo := postgres.NewTodoRepository(pool)
			uc := usecase.NewTodoUseCase(repo, logger)

			count, err := uc.CompleteAllTodos(ctx)
			if err != nil {
				return fmt.Errorf("complete all: %w", err)
			}

			fmt.Printf("Marked %d todos as complete.\n", count)
			return nil
		},
	}

	rootCmd.AddCommand(migrateCmd, listCmd, completeAllCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
