package client

import (
	"context"
	"fmt"

	"github.com/Gitubrr/GoSymGym/gateway/internal/domain"
)

type CollectorClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

var (
	ErrNotFound     = fmt.Errorf("repository not found")
	ErrRateLimit    = fmt.Errorf("rate limit exceeded")
	ErrUnauthorized = fmt.Errorf("unauthorized - check your token")
)
