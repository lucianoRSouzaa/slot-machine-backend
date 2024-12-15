package repository

import "context"

type RefreshTokenRepository interface {
    StoreRefreshToken(ctx context.Context, userID, token string) error
    ValidateRefreshToken(ctx context.Context, userID, token string) (bool, error)
    DeleteRefreshToken(ctx context.Context, userID, token string) error
}
