package jwt

import (
	"errors"
	"fmt"
	"slot-machine/internal/domain/ports"
	"time"

	"github.com/dgrijalva/jwt-go"
)


type JWTManager struct {
	secretKey     string
	accessTokenDuration  time.Duration
    refreshTokenDuration time.Duration
}

func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) ports.JWTManager {
	return &JWTManager{
        secretKey:            secretKey,
        accessTokenDuration:  accessTokenDuration,
        refreshTokenDuration: refreshTokenDuration,
    }
}

func (m *JWTManager) GenerateAccessToken(userID string) (string, error) {
    claims := &ports.JWTClaims{
        UserID:    userID,
        TokenType: ports.TokenTypeAccess,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(m.accessTokenDuration).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
    claims := &ports.JWTClaims{
        UserID:    userID,
        TokenType: ports.TokenTypeRefresh,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(m.refreshTokenDuration).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) VerifyAccessToken(tokenString string) (*ports.JWTClaims, error) {
    claims, err := m.verifyToken(tokenString)
    if err != nil {
        return nil, err
    }
    if claims.TokenType != ports.TokenTypeAccess {
        return nil, errors.New("invalid token type")
    }
    return claims, nil
}

func (m *JWTManager) VerifyRefreshToken(tokenString string) (*ports.JWTClaims, error) {
    claims, err := m.verifyToken(tokenString)
    if err != nil {
        return nil, err
    }
    if claims.TokenType != ports.TokenTypeRefresh {
        return nil, errors.New("invalid token type")
    }
    return claims, nil
}
func (m *JWTManager) verifyToken(tokenString string) (*ports.JWTClaims, error) {
    claims := &ports.JWTClaims{}

    parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(m.secretKey), nil
    })

    if err != nil {
        return nil, err
    }

	fmt.Println(parsedToken)

    if !parsedToken.Valid {
        return nil, errors.New("invalid token")
    }

    return &ports.JWTClaims{
        UserID: claims.UserID,
		TokenType: claims.TokenType,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: claims.ExpiresAt,
            IssuedAt:  claims.IssuedAt,
        },
    }, nil
}
