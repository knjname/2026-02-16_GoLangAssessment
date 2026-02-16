package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/knjname/go-todo-api/internal/config"
	"github.com/knjname/go-todo-api/internal/handler"
	"github.com/knjname/go-todo-api/internal/middleware"
	"github.com/knjname/go-todo-api/internal/usecase"
)

func Run(ctx context.Context, cfg *config.Config, uc *usecase.TodoUseCase, logger *slog.Logger) error {
	mux := http.NewServeMux()

	api := humago.New(mux, huma.DefaultConfig("Todo API", "1.0.0"))

	todoHandler := handler.NewTodoHandler(uc)
	todoHandler.Register(api)

	var h http.Handler = mux
	h = middleware.Logging(logger)(h)
	h = middleware.Recovery(logger)(h)
	h = middleware.RequestID(h)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", slog.Int("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-quit:
		logger.Info("shutdown signal received")
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
