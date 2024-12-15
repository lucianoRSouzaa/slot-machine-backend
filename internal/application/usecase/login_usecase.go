package usecase

import (
	"context"
	"errors"

	"slot-machine/internal/domain/ports"
	"slot-machine/internal/domain/repository"
	"slot-machine/internal/domain/security"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginUseCase struct {
	PlayerRepo repository.PlayerRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	Hasher     security.PasswordHasher
	JWTManager ports.JWTManager
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

func NewLoginUseCase(playerRepo repository.PlayerRepository, refreshRepo repository.RefreshTokenRepository, hasher security.PasswordHasher, jwtManager ports.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		PlayerRepo: playerRepo,
		Hasher:     hasher,
		JWTManager: jwtManager,
		RefreshTokenRepo: refreshRepo,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	player, err := uc.PlayerRepo.GetPlayerByEmail(ctx, req.Email)

	if err != nil {
		if err == repository.ErrPlayerNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = uc.Hasher.CompareHashAndPassword(player.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := uc.JWTManager.GenerateAccessToken(player.ID)
    if err != nil {
        return nil, err
    }

    refreshToken, err := uc.JWTManager.GenerateRefreshToken(player.ID)
    if err != nil {
        return nil, err
    }

	err = uc.RefreshTokenRepo.StoreRefreshToken(ctx, player.ID, refreshToken)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
