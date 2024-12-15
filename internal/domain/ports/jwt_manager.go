package ports

import "github.com/dgrijalva/jwt-go"

type TokenType string

const (
    TokenTypeAccess  TokenType = "access"
    TokenTypeRefresh TokenType = "refresh"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt.StandardClaims
}

type JWTManager interface {
    GenerateAccessToken(userID string) (string, error)
    GenerateRefreshToken(userID string) (string, error)
    VerifyAccessToken(token string) (*JWTClaims, error)
    VerifyRefreshToken(token string) (*JWTClaims, error)
}