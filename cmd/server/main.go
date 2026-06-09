package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omidMorovati/expenseTracker/internal/config"
	"github.com/omidMorovati/expenseTracker/internal/handler"
	"github.com/omidMorovati/expenseTracker/internal/middleware"
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
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	expenseRepo := repository.NewExpenseRepo(pool)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	expenseService := service.NewExpenseService(expenseRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService, logger)
	expenseHandler := handler.NewExpenseHandler(expenseService, logger)

	// Router & Middleware
	r := chi.NewRouter()
	r.Use(slogMiddleware(logger))

	// Public routes
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth([]byte(cfg.JWTSecret)))
		r.Post("/expenses", expenseHandler.Create)
		r.Get("/dashboard", expenseHandler.Dashboard)
	})

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
