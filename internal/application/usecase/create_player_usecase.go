package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"

	"github.com/google/uuid"
)

var (
	ErrPlayerAlreadyExists = errors.New("player already exists")
)

type CreatePlayerUseCase struct {
	PlayerRepo repository.PlayerRepository
}

type CreatePlayerRequest struct {
	Email    string `json:"email"`
	Balance  int    `json:"balance"`
	Password string `json:"password"`
}

type CreatePlayerResponse struct {
	Player model.Player `json:"player"`
}

func NewCreatePlayerUseCase(pr repository.PlayerRepository) *CreatePlayerUseCase {
	return &CreatePlayerUseCase{
		PlayerRepo: pr,
	}
}

func (uc *CreatePlayerUseCase) Execute(ctx context.Context, req *CreatePlayerRequest) (*CreatePlayerResponse, error) {
	playerCreated, err := uc.PlayerRepo.GetPlayerByEmail(ctx, req.Email)

	if playerCreated != nil {
		return nil, ErrPlayerAlreadyExists
	}
	if err != nil && err != repository.ErrPlayerNotFound {
		return nil, err
	}

	player := &model.Player{
		ID:       uuid.New().String(),
		Balance:  req.Balance,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.PlayerRepo.CreatePlayer(ctx, player); err != nil {
		return nil, err
	}

	return &CreatePlayerResponse{
		Player: *player,
	}, nil
}
