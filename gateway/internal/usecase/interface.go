package usecase

import (
	"context"

	"github.com/Gitubrr/GoSymGym/gateway/internal/domain"
)

type RepoUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}
