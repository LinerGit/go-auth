package service

import (
	"context"

	db "github.com/LinerGit/go-auth/internal/repository/db"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		arg db.CreateUserParams,
	) (db.User, error)

	GetUserByUsername(
		ctx context.Context,
		username string,
	) (db.User, error)

	GetUserByID(
		ctx context.Context,
		id int64,
	) (db.User, error)
}

type JWTService interface {
	// FIX 3: added username parameter so it gets embedded in the token
	GenerateAccessToken(
		userID int64,
		username string,
		role string,
	) (string, error)

	GenerateRefreshToken() (string, error)

	HashRefreshToken(token string) string
}

type RefreshRepository interface {
	Create(
		ctx context.Context,
		arg db.CreateRefreshTokenParams,
	) (db.RefreshToken, error)

	GetByHash(
		ctx context.Context,
		hash string,
	) (db.RefreshToken, error)

	DeleteByHash(
		ctx context.Context,
		hash string,
	) error
}
