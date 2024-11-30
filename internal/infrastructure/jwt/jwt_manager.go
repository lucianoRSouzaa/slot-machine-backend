package jwt

import (
	"errors"
	"slot-machine/internal/domain/ports"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) ports.JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (m *JWTManager) Generate(userID string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Verify(token string) (*ports.JWTClaims, error) {
	claims := &JWTClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &ports.JWTClaims{
		UserID: claims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.IssuedAt,
		},
	}, nil
}
