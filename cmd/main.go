package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LinerGit/go-auth/internal/config"
	"github.com/LinerGit/go-auth/internal/database"
	db "github.com/LinerGit/go-auth/internal/repository/db"

	"github.com/LinerGit/go-auth/internal/handler"
	"github.com/LinerGit/go-auth/internal/jwt"
	authmiddleware "github.com/LinerGit/go-auth/internal/middleware"
	"github.com/LinerGit/go-auth/internal/password"
	"github.com/LinerGit/go-auth/internal/repository"
	"github.com/LinerGit/go-auth/internal/server"
	"github.com/LinerGit/go-auth/internal/service"
)

// @title Auth Service API
// @version 1.0
// @description JWT Authentication microservice
// @host localhost:8080
// @BasePath /
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.MustLoad()

	// 1. Postgres pool
	pool, err := database.New(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 2. SQLC layer (ВАЖНО!)
	queries := db.New(pool)

	// repositories
	userRepo := repository.NewUserRepo(queries)
	refreshRepo := repository.NewRefreshRepo(queries)

	// services
	jwtSvc := jwt.New(cfg.JWTSecret, cfg.AccessTTL)
	passwordSvc := password.New(cfg.BcryptCost)

	authSvc := service.NewAuthService(
		userRepo,
		refreshRepo,
		jwtSvc,
		passwordSvc,
		cfg.RefreshTTL,
	)

	// handlers & middleware
	authHandler := handler.NewAuthHandler(authSvc)
	authMw := authmiddleware.NewAuthMiddleware(jwtSvc)

	// router
	r := server.NewRouter(authHandler, authMw)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
}
