package server

import (
	"net/http"
	"time"

	"github.com/LinerGit/go-auth/internal/handler"
	authmiddleware "github.com/LinerGit/go-auth/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	authMw *authmiddleware.AuthMiddleware,
) http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]any{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(authMw.Auth)
			r.Post("/logout", authHandler.Logout)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(authMw.Auth)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("authorized"))
		})
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
