package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret    []byte
	accessTTL time.Duration
}

func New(
	secret string,
	accessTTL time.Duration,
) *JWTService {

	return &JWTService{
		secret:    []byte(secret),
		accessTTL: accessTTL,
	}
}

// FIX 3: added username parameter so chat service can read sender name from token claims
func (s *JWTService) GenerateAccessToken(
	userID int64,
	username string,
	role string,
) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.FormatInt(
				userID,
				10,
			),

			IssuedAt: jwt.NewNumericDate(
				now,
			),

			ExpiresAt: jwt.NewNumericDate(
				now.Add(s.accessTTL),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		s.secret,
	)
}

func (s *JWTService) ValidateAccessToken(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (
			interface{},
			error,
		) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return s.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New(
			"invalid token",
		)
	}

	return claims, nil
}

func (s *JWTService) GenerateRefreshToken() (
	string,
	error,
) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		bytes,
	), nil
}

func (s *JWTService) HashRefreshToken(
	token string,
) string {

	hash := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		hash[:],
	)
}
