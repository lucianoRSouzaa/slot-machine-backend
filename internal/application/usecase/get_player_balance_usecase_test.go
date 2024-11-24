package usecase

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayerBalanceUseCase(t *testing.T) {
	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()

	getPlayerBalanceUC := NewGetPlayerBalanceUseCase(playerRepo)

	ctx := context.Background()

	t.Run("Execute_Success", func(t *testing.T) {
		player := &model.Player{
			ID:      "player1",
			Balance: 1000,
		}
		err := playerRepo.CreatePlayer(ctx, player)
		assert.NoError(t, err, "Expected no error when creating a player")

		req := &GetPlayerBalanceRequest{
			PlayerID: "player1",
		}

		resp, err := getPlayerBalanceUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when getting player balance")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, player.ID, resp.Player.ID, "Player ID should match the request")
		assert.Equal(t, player.Balance, resp.Player.Balance, "Player balance should match the stored balance")
	})

	t.Run("Execute_PlayerNotFound", func(t *testing.T) {
		req := &GetPlayerBalanceRequest{
			PlayerID: "nonexistent_player",
		}

		resp, err := getPlayerBalanceUC.Execute(ctx, req)

		assert.Error(t, err, "Expected an error when getting balance of a non-existent player")
		assert.Equal(t, repository.ErrPlayerNotFound, err, "Expected ErrPlayerNotFound error")
		assert.Nil(t, resp, "Expected no response when there is an error")
	})
}
