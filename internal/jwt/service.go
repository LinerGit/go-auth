package jwt

type Service interface {
	GenerateAccessToken(
		userID int64,
		role string,
	) (string, error)

	ValidateAccessToken(
		tokenString string,
	) (*Claims, error)

	GenerateRefreshToken() (string, error)

	HashRefreshToken(
		token string,
	) string
}
