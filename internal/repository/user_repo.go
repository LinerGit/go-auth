package repository

import (
	"context"

	db "github.com/LinerGit/go-auth/internal/repository/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByID(ctx context.Context, id int64) (db.User, error)
}

type userRepo struct {
	q *db.Queries
}

func NewUserRepo(q *db.Queries) UserRepository {
	return &userRepo{q: q}
}

func (r *userRepo) CreateUser(
	ctx context.Context,
	arg db.CreateUserParams,
) (db.User, error) {

	return r.q.CreateUser(ctx, arg)
}

func (r *userRepo) GetUserByUsername(
	ctx context.Context,
	username string,
) (db.User, error) {

	return r.q.GetUserByUsername(ctx, username)
}

func (r *userRepo) GetUserByID(
	ctx context.Context,
	id int64,
) (db.User, error) {

	return r.q.GetUserByID(ctx, id)
}
