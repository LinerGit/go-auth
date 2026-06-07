package service

import (
	"context"
	"errors"
	"time"

	db "github.com/LinerGit/go-auth/internal/repository/db"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	users         UserRepository
	refreshTokens RefreshRepository

	jwt      JWTService
	password PasswordService

	refreshTTL time.Duration
}

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

func NewAuthService(
	users UserRepository,
	refreshTokens RefreshRepository,
	jwt JWTService,
	password PasswordService,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		jwt:           jwt,
		password:      password,
		refreshTTL:    refreshTTL,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	username string,
	password string,
) (*AuthTokens, error) {

	_, err := s.users.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := s.password.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := s.users.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		Role:         "USER",
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// fix: was ignoring error from GenerateRefreshToken
	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// fix: was calling SaveRefreshToken (which calls Create internally)
	// AND then calling refreshTokens.Create again — double insert removed
	if err = s.SaveRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshTTL)); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (*AuthTokens, error) {

	user, err := s.users.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err = s.password.Compare(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// fix: was ignoring error from GenerateRefreshToken
	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// fix: same double insert bug as Register — removed
	if err = s.SaveRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshTTL)); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh validates the old refresh token, rotates it, and returns a new token pair.
func (s *AuthService) Refresh(
	ctx context.Context,
	rawRefreshToken string,
) (*AuthTokens, error) {

	storedToken, err := s.ValidateRefreshToken(ctx, rawRefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// fetch user to get role
	user, err := s.users.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	// rotate: delete old token
	if err = s.refreshTokens.DeleteByHash(ctx, storedToken.TokenHash); err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err = s.SaveRefreshToken(ctx, user.ID, newRefreshToken, time.Now().Add(s.refreshTTL)); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	rawRefreshToken string,
) error {
	hash := s.jwt.HashRefreshToken(rawRefreshToken)
	return s.refreshTokens.DeleteByHash(ctx, hash)
}

func (s *AuthService) SaveRefreshToken(
	ctx context.Context,
	userID int64,
	rawToken string,
	expiresAt time.Time,
) error {
	hash := s.jwt.HashRefreshToken(rawToken)
	_, err := s.refreshTokens.Create(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: expiresAt,
	})
	return err
}

func (s *AuthService) ValidateRefreshToken(
	ctx context.Context,
	rawToken string,
) (db.RefreshToken, error) {
	hash := s.jwt.HashRefreshToken(rawToken)

	token, err := s.refreshTokens.GetByHash(ctx, hash)
	if err != nil {
		return db.RefreshToken{}, err
	}

	if time.Now().After(token.ExpiresAt) {
		return db.RefreshToken{}, ErrInvalidToken
	}

	return token, nil
}
