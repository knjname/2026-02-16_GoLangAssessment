package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/knjname/go-todo-api/internal/di"
	"github.com/knjname/go-todo-api/internal/server"
)

func main() {
	ctx := context.Background()

	components, err := di.InitializeAPI(ctx)
	if err != nil {
		slog.Error("failed to initialize", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer components.Pool.Close()

	if err := server.Run(ctx, components.Config, components.UseCase, components.Logger); err != nil {
		components.Logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
