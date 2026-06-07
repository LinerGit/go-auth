package server

import (
	"net/http"
	"time"

	"github.com/LinerGit/go-auth/internal/handler"
	authmiddleware "github.com/LinerGit/go-auth/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	authMw *authmiddleware.AuthMiddleware,
) http.Handler {

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(authMw.Auth)
			r.Post("/logout", authHandler.Logout)
		})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMw.Auth)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("authorized"))
		})
	})

	return r
}
