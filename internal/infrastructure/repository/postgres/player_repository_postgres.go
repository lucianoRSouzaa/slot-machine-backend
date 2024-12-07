package repository_postgres

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPlayerRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPlayerRepository(pool *pgxpool.Pool) repository.PlayerRepository {
	return &PostgresPlayerRepository{
		pool: pool,
	}
}

func (r *PostgresPlayerRepository) CreatePlayer(ctx context.Context, player *model.Player) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO players (id, balance, email, password, role)
		VALUES ($1, $2, $3, $4, $5)`,
		player.ID, player.Balance, player.Email, player.Password, player.Role)
	return err
}

func (r *PostgresPlayerRepository) GetPlayer(ctx context.Context, id string) (*model.Player, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, balance, email, password, role
		FROM players
		WHERE id = $1`, id)
	player := &model.Player{}
	err := row.Scan(&player.ID, &player.Balance, &player.Email, &player.Password, &player.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrPlayerNotFound
		}
		return nil, err
	}
	return player, nil
}

func (r *PostgresPlayerRepository) GetPlayerByEmail(ctx context.Context, email string) (*model.Player, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, balance, email, password, role
		FROM players
		WHERE email = $1`, email)
	player := &model.Player{}
	err := row.Scan(&player.ID, &player.Balance, &player.Email, &player.Password, &player.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrPlayerNotFound
		}
		return nil, err
	}
	return player, nil
}

func (r *PostgresPlayerRepository) UpdatePlayer(ctx context.Context, player *model.Player) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE players
		SET balance = $1, email = $2, password = $3, role = $4
		WHERE id = $5`,
		player.Balance, player.Email, player.Password, player.Role, player.ID)
	return err
}

func (r *PostgresPlayerRepository) ListPlayers(ctx context.Context) ([]*model.Player, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, balance, email, password, role
		FROM players`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*model.Player
	for rows.Next() {
		player := &model.Player{}
		err := rows.Scan(&player.ID, &player.Balance, &player.Email, &player.Password, &player.Role)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}
