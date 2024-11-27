package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"slot-machine/internal/domain/security"

	"github.com/google/uuid"
)

var (
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrValidate            = errors.New("balance, email, and password must be provided")
)

type CreatePlayerUseCase struct {
	PlayerRepo     repository.PlayerRepository
	PasswordHasher security.PasswordHasher
}

type CreatePlayerRequest struct {
	Email    string `json:"email"`
	Balance  int    `json:"balance"`
	Password string `json:"password"`
}

type CreatePlayerResponse struct {
	Player model.Player `json:"player"`
}

func NewCreatePlayerUseCase(repo repository.PlayerRepository, hasher security.PasswordHasher) *CreatePlayerUseCase {
	return &CreatePlayerUseCase{
		PlayerRepo:     repo,
		PasswordHasher: hasher,
	}
}

func (uc *CreatePlayerUseCase) Execute(ctx context.Context, req *CreatePlayerRequest) (*CreatePlayerResponse, error) {

	if req.Email == "" || req.Password == "" {
		return nil, ErrValidate
	}

	playerCreated, err := uc.PlayerRepo.GetPlayerByEmail(ctx, req.Email)

	if playerCreated != nil {
		return nil, ErrPlayerAlreadyExists
	}
	if err != nil && err != repository.ErrPlayerNotFound {
		return nil, err
	}

	passwordHashed, err := uc.PasswordHasher.Hash(req.Password)

	if err != nil {
		return nil, err
	}

	player := &model.Player{
		ID:       uuid.New().String(),
		Balance:  req.Balance,
		Email:    req.Email,
		Password: passwordHashed,
	}

	if err := uc.PlayerRepo.CreatePlayer(ctx, player); err != nil {
		return nil, err
	}

	player.Password = ""

	return &CreatePlayerResponse{
		Player: *player,
	}, nil
}
