package repository

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, player *model.Player) error
	GetPlayer(ctx context.Context, id string) (*model.Player, error)
	GetPlayerByEmail(ctx context.Context, email string) (*model.Player, error)
	UpdatePlayer(ctx context.Context, player *model.Player) error
	ListPlayers(ctx context.Context) ([]*model.Player, error)
}
