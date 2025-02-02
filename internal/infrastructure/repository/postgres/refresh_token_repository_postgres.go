package repository_postgres

import (
	"context"

	"slot-machine/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRefreshTokenRepository(pool *pgxpool.Pool) repository.RefreshTokenRepository {
	return &PostgresRefreshTokenRepository{pool: pool}
}

func (r *PostgresRefreshTokenRepository) StoreRefreshToken(ctx context.Context, userID string, token string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, token)
	return err
}

func (r *PostgresRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM refresh_tokens
		WHERE user_id = $1 AND token = $2
	`, userID, token)
	return err
}

func (r *PostgresRefreshTokenRepository) ValidateRefreshToken(ctx context.Context, userID string, token string) (bool, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT token
		FROM refresh_tokens
		WHERE user_id = $1 AND token = $2
	`, userID, token)

	var retrievedToken string
	err := row.Scan(&retrievedToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
