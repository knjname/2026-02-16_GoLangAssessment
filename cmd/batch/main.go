package main

import (
	"context"
	"fmt"
	"os"

	"github.com/knjname/go-todo-api/internal/di"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func main() {
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
			ctx := context.Background()
			components, err := di.InitializeBatch(ctx)
			if err != nil {
				return fmt.Errorf("initialize: %w", err)
			}
			defer components.Pool.Close()
			defer func() { _ = components.DB.Close() }()

			if err := goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("set dialect: %w", err)
			}
			if err := goose.Up(components.DB, "migrations"); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}
			components.Logger.Info("migrations applied successfully")
			return nil
		},
	}

	migrateDownCmd := &cobra.Command{
		Use:   "down",
		Short: "Roll back the last migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			components, err := di.InitializeBatch(ctx)
			if err != nil {
				return fmt.Errorf("initialize: %w", err)
			}
			defer components.Pool.Close()
			defer func() { _ = components.DB.Close() }()

			if err := goose.SetDialect("postgres"); err != nil {
				return fmt.Errorf("set dialect: %w", err)
			}
			if err := goose.Down(components.DB, "migrations"); err != nil {
				return fmt.Errorf("migrate down: %w", err)
			}
			components.Logger.Info("migration rolled back successfully")
			return nil
		},
	}

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			components, err := di.InitializeBatch(ctx)
			if err != nil {
				return fmt.Errorf("initialize: %w", err)
			}
			defer components.Pool.Close()
			defer func() { _ = components.DB.Close() }()

			todos, err := components.UseCase.ListTodos(ctx)
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
			components, err := di.InitializeBatch(ctx)
			if err != nil {
				return fmt.Errorf("initialize: %w", err)
			}
			defer components.Pool.Close()
			defer func() { _ = components.DB.Close() }()

			count, err := components.UseCase.CompleteAllTodos(ctx)
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
