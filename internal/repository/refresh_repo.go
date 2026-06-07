package repository

import (
	"context"

	db "github.com/LinerGit/go-auth/internal/repository/db"
)

type RefreshRepo struct {
	q *db.Queries
}

func NewRefreshRepo(q *db.Queries) *RefreshRepo {
	return &RefreshRepo{q: q}
}

func (r *RefreshRepo) Create(
	ctx context.Context,
	arg db.CreateRefreshTokenParams,
) (db.RefreshToken, error) {

	return r.q.CreateRefreshToken(ctx, arg)
}

func (r *RefreshRepo) GetByHash(
	ctx context.Context,
	hash string,
) (db.RefreshToken, error) {

	return r.q.GetRefreshToken(ctx, hash)
}

func (r *RefreshRepo) DeleteByHash(
	ctx context.Context,
	hash string,
) error {

	return r.q.DeleteRefreshToken(ctx, hash)
}
