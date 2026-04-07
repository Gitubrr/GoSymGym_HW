package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Gitubrr/GoSymGym/collector/internal/domain"
	"github.com/Gitubrr/GoSymGym/collector/internal/repository"
)

var (
	ErrNotFound     = errors.New("repository not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrRateLimit    = errors.New("rate limit exceeded")
	ErrUnauthorized = errors.New("unauthorized")
)

type GitHubUseCase struct {
	githubRepo repository.GitHubRepository
}

func NewGitHubUseCase(githubRepo repository.GitHubRepository) *GitHubUseCase {
	return &GitHubUseCase{
		githubRepo: githubRepo,
	}
}

func (uc *GitHubUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	log.Printf("GetRepository called with owner=%s, repo=%s", owner, repo)

	if owner == "" {
		return nil, fmt.Errorf("%w: owner is required", ErrInvalidInput)
	}
	if repo == "" {
		return nil, fmt.Errorf("%w: repo is required", ErrInvalidInput)
	}

	repoData, err := uc.githubRepo.GetRepository(ctx, owner, repo)
	if err != nil {
		log.Printf("Repository error: %v", err)

		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrNotFound
		case errors.Is(err, repository.ErrRateLimit):
			return nil, ErrRateLimit
		case errors.Is(err, repository.ErrUnauthorized):
			return nil, ErrUnauthorized
		default:
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}
	}

	return repoData, nil
}
