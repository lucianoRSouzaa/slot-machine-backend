package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
)

var (
	ErrPlayerAlreadyExists = errors.New("player already exists")
)

type CreatePlayerUseCase struct {
	PlayerRepo repository.PlayerRepository
}

type CreatePlayerRequest struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
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
	_, err := uc.PlayerRepo.GetPlayer(ctx, req.ID)
	if err == nil {
		return nil, ErrPlayerAlreadyExists
	}
	if err != repository.ErrPlayerNotFound {
		return nil, err
	}

	player := &model.Player{
		ID:      req.ID,
		Balance: req.Balance,
	}

	if err := uc.PlayerRepo.CreatePlayer(ctx, player); err != nil {
		return nil, err
	}

	return &CreatePlayerResponse{
		Player: *player,
	}, nil
}
