package security

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	Cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	return &BcryptPasswordHasher{Cost: cost}
}

func (b *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), b.Cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (b *BcryptPasswordHasher) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
