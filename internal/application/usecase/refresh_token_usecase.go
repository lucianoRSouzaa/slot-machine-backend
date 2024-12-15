package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/ports"
	"slot-machine/internal/domain/repository"
)

var (
    ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUseCase struct {
    JWTManager        ports.JWTManager
    RefreshTokenRepo  repository.RefreshTokenRepository
}

func NewRefreshTokenUseCase(jwtManager ports.JWTManager, refreshTokenRepo repository.RefreshTokenRepository) *RefreshTokenUseCase {
    return &RefreshTokenUseCase{
        JWTManager:       jwtManager,
        RefreshTokenRepo: refreshTokenRepo,
    }
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
    claims, err := uc.JWTManager.VerifyRefreshToken(req.RefreshToken)
    if err != nil {
        return nil, ErrInvalidRefreshToken
    }

    valid, err := uc.RefreshTokenRepo.ValidateRefreshToken(ctx, claims.UserID, req.RefreshToken)
    if err != nil || !valid {
        return nil, ErrInvalidRefreshToken
    }

    newAccessToken, err := uc.JWTManager.GenerateAccessToken(claims.UserID)
    if err != nil {
        return nil, err
    }

    newRefreshToken, err := uc.JWTManager.GenerateRefreshToken(claims.UserID)
    if err != nil {
        return nil, err
    }

    _ = uc.RefreshTokenRepo.DeleteRefreshToken(ctx, claims.UserID, req.RefreshToken)

    err = uc.RefreshTokenRepo.StoreRefreshToken(ctx, claims.UserID, newRefreshToken)
    if err != nil {
        return nil, err
    }

    return &RefreshTokenResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
    }, nil
}
