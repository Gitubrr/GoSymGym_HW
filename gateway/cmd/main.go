package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Gitubrr/GoSymGym/gateway/internal/client"
	"github.com/Gitubrr/GoSymGym/gateway/internal/delivery/rest"
	"github.com/Gitubrr/GoSymGym/gateway/internal/usecase"
)

type config struct {
	httpPort      string
	collectorAddr string
}

func loadConfig() config {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}

	return config{
		httpPort:      httpPort,
		collectorAddr: collectorAddr,
	}
}

func main() {
	cfg := loadConfig()

	// Создаем gRPC клиент к Collector
	collectorClient, err := client.NewCollectorClient(cfg.collectorAddr)
	if err != nil {
		log.Fatalf("Failed to create collector client: %v", err)
	}
	defer collectorClient.Close()

	// Создаем usecase
	repoUseCase := usecase.NewRepoUseCase(collectorClient)

	// Создаем HTTP handler
	handler := rest.NewHandler(repoUseCase)

	// Настраиваем роутер
	router := rest.SetupRouter(handler)

	// Запускаем сервер
	log.Printf("Starting Gateway HTTP server on port %s", cfg.httpPort)
	log.Printf("Connecting to Collector at %s", cfg.collectorAddr)

	if err := http.ListenAndServe(":"+cfg.httpPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
