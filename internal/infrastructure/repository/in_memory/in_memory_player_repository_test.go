package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryPlayerRepository(t *testing.T) {
	repo := NewInMemoryPlayerRepository()

	ctx := context.Background()

	t.Run("CreatePlayer_Success", func(t *testing.T) {
		player := &model.Player{
			ID:      "player1",
			Balance: 1000,
		}

		err := repo.CreatePlayer(ctx, player)
		assert.NoError(t, err, "Expected no error on creating player")

		retrievedPlayer, err := repo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Expected no error on retrieving existing player")
		assert.Equal(t, player, retrievedPlayer, "Retrieved player should match the created player")
	})

	t.Run("CreatePlayer_DuplicateID", func(t *testing.T) {
		player := &model.Player{
			ID:      "player1",
			Balance: 2000,
		}

		err := repo.CreatePlayer(ctx, player)
		assert.NoError(t, err, "Expected no error on creating player with duplicate ID")

		retrievedPlayer, err := repo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Expected no error on retrieving existing player")
		assert.Equal(t, player, retrievedPlayer, "Retrieved player should match the new player data")
	})

	t.Run("GetPlayer_NotFound", func(t *testing.T) {
		_, err := repo.GetPlayer(ctx, "nonexistent_player")
		assert.Error(t, err, "Expected error when retrieving non-existent player")
		assert.Equal(t, repository.ErrPlayerNotFound, err, "Expected ErrPlayerNotFound error")
	})

	t.Run("UpdatePlayer_Success", func(t *testing.T) {
		player := &model.Player{
			ID:      "player1",
			Balance: 1500,
		}

		err := repo.UpdatePlayer(ctx, player)
		assert.NoError(t, err, "Expected no error on updating existing player")

		retrievedPlayer, err := repo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Expected no error on retrieving existing player")
		assert.Equal(t, player, retrievedPlayer, "Retrieved player should reflect the updated data")
	})

	t.Run("UpdatePlayer_NotFound", func(t *testing.T) {
		player := &model.Player{
			ID:      "nonexistent_player",
			Balance: 500,
		}

		err := repo.UpdatePlayer(ctx, player)
		assert.Error(t, err, "Expected error on updating non-existent player")
		assert.Equal(t, repository.ErrPlayerNotFound, err, "Expected ErrPlayerNotFound error")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100
		playerIDs := make([]string, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			playerID := "concurrent_player_" + strconv.Itoa(i)
			playerIDs[i] = playerID
			go func(id string) {
				defer wg.Done()
				player := &model.Player{
					ID:      id,
					Balance: 100,
				}
				err := repo.CreatePlayer(ctx, player)
				assert.NoError(t, err, "Expected no error on concurrent player creation")
			}(playerID)
		}

		wg.Wait()

		for _, id := range playerIDs {
			player, err := repo.GetPlayer(ctx, id)
			assert.NoError(t, err, "Expected no error on retrieving concurrently created player")
			assert.Equal(t, 100, player.Balance, "Expected player balance to be 100")
		}
	})
}
