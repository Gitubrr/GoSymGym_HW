package main

import (
	"log"
	"os"
	"strconv"

	grpc "github.com/Gitubrr/GoSymGym/collector/internal/delivery/grpc"
	"github.com/Gitubrr/GoSymGym/collector/internal/repository"
	"github.com/Gitubrr/GoSymGym/collector/internal/usecase"
)

type serverConfig struct {
	gRPCPort string
	token    string
	timeout  int
}

func loadServerConfig() serverConfig {
	gRPCPort := os.Getenv("GRPC_PORT")
	if gRPCPort == "" {
		gRPCPort = "50051"
		log.Printf("GRPC_PORT not set, using default: %s", gRPCPort)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("GITHUB_TOKEN not set, requests will be unauthenticated")
	}

	timeout := 10
	if timeoutStr := os.Getenv("REQUEST_TIMEOUT"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		} else {
			log.Printf("Invalid REQUEST_TIMEOUT: %s", timeoutStr)
		}
	}

	return serverConfig{
		gRPCPort: gRPCPort,
		token:    token,
		timeout:  timeout,
	}
}

func main() {

	config := loadServerConfig()

	client := repository.NewGitHubClient(
		config.token,
		config.timeout,
	)
	getRepoUseCase := usecase.NewGitHubUseCase(client)
	grpcHandler := grpc.NewHandler(getRepoUseCase)

	log.Printf("Starting Collector gRPC server on port %s", config.gRPCPort)
	if err := grpc.RunServer(config.gRPCPort, grpcHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
