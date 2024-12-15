package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/repository"
	"sync"
)

type InMemoryRefreshTokenRepository struct {
	tokens map[string][]string
	mu     sync.Mutex
}

func NewInMemoryRefreshTokenRepository() repository.RefreshTokenRepository {
	return &InMemoryRefreshTokenRepository{
		tokens: make(map[string][]string),
	}
}

func (r *InMemoryRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, token string) error {
	r.mu.Lock()
    defer r.mu.Unlock()

    if tokens, ok := r.tokens[userID]; ok {
        newTokens := make([]string, 0, len(tokens))
        for _, t := range tokens {
            if t != token {
                newTokens = append(newTokens, t)
            }
        }
        r.tokens[userID] = newTokens
    }
    return nil
}

func (r *InMemoryRefreshTokenRepository) StoreRefreshToken(ctx context.Context, userID string, token string) error {
	r.mu.Lock()
    defer r.mu.Unlock()

    r.tokens[userID] = append(r.tokens[userID], token)
    return nil
}

func (r *InMemoryRefreshTokenRepository) ValidateRefreshToken(ctx context.Context, userID string, token string) (bool, error) {
	r.mu.Lock()
    defer r.mu.Unlock()

    if tokens, ok := r.tokens[userID]; ok {
        for _, t := range tokens {
            if t == token {
                return true, nil
            }
        }
    }
    return false, nil
}

