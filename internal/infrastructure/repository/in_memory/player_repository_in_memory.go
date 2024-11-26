package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"sync"
)

type InMemoryPlayerRepository struct {
	players map[string]*model.Player
	mu      sync.RWMutex
}

func NewInMemoryPlayerRepository() repository.PlayerRepository {
	return &InMemoryPlayerRepository{
		players: make(map[string]*model.Player),
	}
}

func (r *InMemoryPlayerRepository) CreatePlayer(ctx context.Context, player *model.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.players[player.ID] = player
	return nil
}

func (r *InMemoryPlayerRepository) GetPlayer(ctx context.Context, id string) (*model.Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	player, exists := r.players[id]
	if !exists {
		return nil, repository.ErrPlayerNotFound
	}
	return player, nil
}

func (r *InMemoryPlayerRepository) UpdatePlayer(ctx context.Context, player *model.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.players[player.ID]
	if !exists {
		return repository.ErrPlayerNotFound
	}
	r.players[player.ID] = player
	return nil
}

func (r *InMemoryPlayerRepository) GetPlayerByEmail(ctx context.Context, email string) (*model.Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, player := range r.players {
		if player.Email == email {
			return player, nil
		}
	}

	return nil, repository.ErrPlayerNotFound
}
