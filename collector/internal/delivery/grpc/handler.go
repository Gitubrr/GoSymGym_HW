package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Gitubrr/GoSymGym/api/collector"
	"github.com/Gitubrr/GoSymGym/collector/internal/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	repoUseCase usecase.RepoUseCase
}

func NewHandler(repoUseCase usecase.RepoUseCase) *Handler {
	return &Handler{
		repoUseCase: repoUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {

	owner := req.GetOwner()
	repo := req.GetRepo()

	repoData, err := h.repoUseCase.GetRepository(ctx, owner, repo)
	if err != nil {

		switch {
		case errors.Is(err, usecase.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, usecase.ErrNotFound):
			return nil, status.Error(codes.NotFound, "repository not found")
		case errors.Is(err, usecase.ErrRateLimit):
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded, please try again later")
		case errors.Is(err, usecase.ErrUnauthorized):
			return nil, status.Error(codes.Unauthenticated, "unauthorized, check your GitHub token")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.RepoResponse{
		Name:        repoData.Name,
		FullName:    repoData.FullName,
		Description: repoData.Description,
		Stars:       int32(repoData.Stars),
		Forks:       int32(repoData.Forks),
		CreatedAt:   timestamppb.New(repoData.CreatedAt).String(),
	}, nil
}
