package jwt

type Service interface {
	// FIX 3: added username parameter to match updated GenerateAccessToken signature
	GenerateAccessToken(
		userID int64,
		username string,
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
