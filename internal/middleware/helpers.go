package middleware

import (
	"context"
	"net/http"
)

func GetUserID(ctx context.Context) (int64, bool) {

	val := ctx.Value(UserIDKey)

	id, ok := val.(int64)

	return id, ok
}

func GetRole(ctx context.Context) (string, bool) {

	val := ctx.Value(RoleKey)

	role, ok := val.(string)

	return role, ok
}

func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			userRole, ok := GetRole(ctx)

			if !ok {
				http.Error(w, "no role", http.StatusForbidden)
				return
			}

			if userRole != role {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
