package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/LinerGit/go-auth/internal/jwt"
)

type AuthMiddleware struct {
	jwt jwt.Service
}

func NewAuthMiddleware(jwtSvc jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwtSvc}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.jwt.ValidateAccessToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		ctx = contextWithUserID(ctx, claims.UserID)
		ctx = contextWithRole(ctx, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func contextWithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, RoleKey, role)
}
