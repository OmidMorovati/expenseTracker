package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omidMorovati/expenseTracker/internal/config"
	"github.com/omidMorovati/expenseTracker/internal/handler"
	"github.com/omidMorovati/expenseTracker/internal/repository"
	"github.com/omidMorovati/expenseTracker/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
	}
	defer pool.Close()

	repo := repository.NewExpenseRepo(pool)
	svc := service.NewExpenseService(repo)
	h := handler.NewExpenseHandler(svc, logger)

	r := chi.NewRouter()
	r.Use(slogMiddleware(logger))
	r.Post("/expenses", h.Create)
	r.Get("/dashboard", h.Dashboard)
	// Add report endpoints...

	srv := &http.Server{Addr: cfg.Port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server crashed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		return
	}
}

func slogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info(r.Method+" "+r.URL.Path,
				"duration", time.Since(start).String(),
				"status", http.StatusText(200),
			)
		})
	}
}
