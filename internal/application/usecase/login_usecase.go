package usecase

import (
	"context"
	"errors"
	"fmt"

	"slot-machine/internal/domain/ports"
	"slot-machine/internal/domain/repository"
	"slot-machine/internal/domain/security"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginUseCase struct {
	PlayerRepo repository.PlayerRepository
	Hasher     security.PasswordHasher
	JWTManager ports.JWTManager
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginUseCase(repo repository.PlayerRepository, hasher security.PasswordHasher, jwtManager ports.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		PlayerRepo: repo,
		Hasher:     hasher,
		JWTManager: jwtManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	players, err := uc.PlayerRepo.ListPlayers(ctx)
	player, err := uc.PlayerRepo.GetPlayerByEmail(ctx, req.Email)

	fmt.Println(players)
	fmt.Println(player)

	if err != nil {
		if err == repository.ErrPlayerNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	fmt.Println(player.Password)
	fmt.Println(req.Password)

	err = uc.Hasher.CompareHashAndPassword(player.Password, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := uc.JWTManager.Generate(player.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
	}, nil
}
