package usecase

import (
	"context"
	"slot-machine/internal/domain/model"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"slot-machine/internal/infrastructure/security"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreatePlayerUseCase(t *testing.T) {
	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()
	hasher := security.NewBcryptPasswordHasher(bcrypt.DefaultCost)

	createPlayerUC := NewCreatePlayerUseCase(playerRepo, hasher)

	ctx := context.Background()

	t.Run("Execute_Success", func(t *testing.T) {
		req := &CreatePlayerRequest{
			Balance:  1000,
			Email:    "email@email.co",
			Password: "password",
		}

		resp, err := createPlayerUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when creating a new player")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, req.Email, resp.Player.Email, "Player Email should match the request")
		assert.Equal(t, req.Balance, resp.Player.Balance, "Player balance should match the request")

		storedPlayer, err := playerRepo.GetPlayer(ctx, resp.Player.ID)
		assert.NoError(t, err, "Expected no error when retrieving the created player")
		assert.Equal(t, resp.Player, *storedPlayer, "Stored player should match the response")
	})

	t.Run("Execute_PlayerAlreadyExists", func(t *testing.T) {
		initialPlayer := &model.Player{
			ID:       "player2",
			Balance:  500,
			Email:    "email@email.co",
			Password: "password",
		}
		err := playerRepo.CreatePlayer(ctx, initialPlayer)
		assert.NoError(t, err, "Expected no error when initially creating a player")

		req := &CreatePlayerRequest{
			Balance:  1500,
			Email:    "email@email.co",
			Password: "password",
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
