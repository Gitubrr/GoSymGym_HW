package client

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("Connecting to collector at %s", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}
	log.Printf("Connected successfully")
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
	log.Printf("Calling GetRepository with owner=%s, repo=%s", owner, repo)
	req := &pb.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := c.client.GetRepository(ctx, req)
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}

	// Парсим время
	log.Printf("gRPC call succeeded, response: %+v", resp)
	createdAt := resp.CreatedAt.AsTime()
	if err != nil {
		log.Printf("Failed to parse time: %v", err)
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
