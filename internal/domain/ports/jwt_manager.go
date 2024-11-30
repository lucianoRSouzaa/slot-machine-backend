package ports

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type JWTManager interface {
	Generate(userID string) (string, error)
	Verify(token string) (*JWTClaims, error)
}
