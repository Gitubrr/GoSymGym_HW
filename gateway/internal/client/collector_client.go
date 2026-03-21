package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Gitubrr/GoSymGym/api/collector"
	"github.com/Gitubrr/GoSymGym/gateway/internal/domain"
)

type collectorClient struct {
	client pb.CollectorServiceClient
	conn   *grpc.ClientConn
}

func NewCollectorClient(addr string) (*collectorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	return &collectorClient{
		client: pb.NewCollectorServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *collectorClient) Close() error {
	return c.conn.Close()
}

// parseTime преобразует строку в time.Time
func parseTime(timeStr string) (time.Time, error) {
	// GitHub API использует RFC3339 формат
	return time.Parse(time.RFC3339, timeStr)
}

func (c *collectorClient) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	req := &pb.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.client.GetRepository(ctx, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}

	// Парсим время
	createdAt, err := parseTime(resp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &domain.Repository{
		Name:        resp.Name,
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   createdAt,
	}, nil
}
