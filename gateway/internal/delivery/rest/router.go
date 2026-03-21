package rest

import (
	"net/http"
)

func SetupRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("GET /api/v1/repos/{owner}/{repo}", handler.GetRepository)

	// Swagger UI (если есть статические файлы)
	mux.HandleFunc("GET /swagger/", serveSwagger)

	return mux
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.yaml")
}
