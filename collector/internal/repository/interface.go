package repository

import (
	"context"
	"fmt"

	"github.com/Gitubrr/GoSymGym/collector/internal/domain"
)

type GitHubRepository interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

var (
	ErrNotFound     = fmt.Errorf("repository not found")
	ErrRateLimit    = fmt.Errorf("rate limit exceeded")
	ErrUnauthorized = fmt.Errorf("unauthorized - check your token")
)
