package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"` // FIX 3: was missing — chat service needs this to display sender name
	Role     string `json:"role"`

	jwt.RegisteredClaims
}
