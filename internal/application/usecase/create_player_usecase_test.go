package usecase

import (
	"context"
	"slot-machine/internal/domain/model"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlayerUseCase(t *testing.T) {
	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()

	createPlayerUC := NewCreatePlayerUseCase(playerRepo)

	ctx := context.Background()

	t.Run("Execute_Success", func(t *testing.T) {
		req := &CreatePlayerRequest{
			ID:      "player1",
			Balance: 1000,
		}

		resp, err := createPlayerUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when creating a new player")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, req.ID, resp.Player.ID, "Player ID should match the request")
		assert.Equal(t, req.Balance, resp.Player.Balance, "Player balance should match the request")

		storedPlayer, err := playerRepo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Expected no error when retrieving the created player")
		assert.Equal(t, resp.Player, *storedPlayer, "Stored player should match the response")
	})

	t.Run("Execute_PlayerAlreadyExists", func(t *testing.T) {
		initialPlayer := &model.Player{
			ID:      "player2",
			Balance: 500,
		}
		err := playerRepo.CreatePlayer(ctx, initialPlayer)
		assert.NoError(t, err, "Expected no error when initially creating a player")

		req := &CreatePlayerRequest{
			ID:      "player2",
			Balance: 1500,
		}

		resp, err := createPlayerUC.Execute(ctx, req)

		assert.Error(t, err, "Expected an error when creating a player that already exists")
		assert.Equal(t, ErrPlayerAlreadyExists, err, "Expected ErrPlayerAlreadyExists error")
		assert.Nil(t, resp, "Expected no response when there is an error")

		storedPlayer, err := playerRepo.GetPlayer(ctx, "player2")
		assert.NoError(t, err, "Expected no error when retrieving the existing player")
		assert.Equal(t, initialPlayer.Balance, storedPlayer.Balance, "Player balance should remain unchanged")
	})
}
