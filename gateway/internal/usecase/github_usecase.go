package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gitubrr/GoSymGym/gateway/internal/client"
	"github.com/Gitubrr/GoSymGym/gateway/internal/domain"
)

var (
	ErrNotFound     = errors.New("repository not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnavailable  = errors.New("collector service unavailable")
)

type repoUseCase struct {
	collectorClient client.CollectorClient
}

func NewRepoUseCase(collectorClient client.CollectorClient) RepoUseCase {
	return &repoUseCase{
		collectorClient: collectorClient,
	}
}

func (uc *repoUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	// Валидация
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("%w: owner and repo are required", ErrInvalidInput)
	}

	// Вызов Collector через gRPC
	repoData, err := uc.collectorClient.GetRepository(ctx, owner, repo)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}

	return repoData, nil
}
