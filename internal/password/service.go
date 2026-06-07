package password

type Service interface {
	Hash(password string) (string, error)

	Compare(
		hashedPassword string,
		password string,
	) error
}
