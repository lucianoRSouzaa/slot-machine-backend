package usecase

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
)

type GetPlayerBalanceUseCase struct {
	PlayerRepo repository.PlayerRepository
}

type GetPlayerBalanceRequest struct {
	PlayerID string `json:"player_id"`
}

type GetPlayerBalanceResponse struct {
	Player model.Player `json:"player"`
}

func NewGetPlayerBalanceUseCase(pr repository.PlayerRepository) *GetPlayerBalanceUseCase {
	return &GetPlayerBalanceUseCase{
		PlayerRepo: pr,
	}
}

func (uc *GetPlayerBalanceUseCase) Execute(ctx context.Context, req *GetPlayerBalanceRequest) (*GetPlayerBalanceResponse, error) {
	player, err := uc.PlayerRepo.GetPlayer(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	return &GetPlayerBalanceResponse{
		Player: *player,
	}, nil
}
