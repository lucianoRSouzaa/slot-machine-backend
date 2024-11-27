package security

type PasswordHasher interface {
	Hash(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}
