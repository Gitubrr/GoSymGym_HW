package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Gitubrr/GoSymGym/gateway/internal/usecase"
)

type Handler struct {
	repoUseCase usecase.RepoUseCase
}

func NewHandler(repoUseCase usecase.RepoUseCase) *Handler {
	return &Handler{
		repoUseCase: repoUseCase,
	}
}

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	// Вызываем usecase
	repoData, err := h.repoUseCase.GetRepository(r.Context(), owner, repo)
	if err != nil {
		// Маппинг ошибок в HTTP статусы
		switch {
		case errors.Is(err, usecase.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, usecase.ErrNotFound):
			http.Error(w, "repository not found", http.StatusNotFound)
		case errors.Is(err, usecase.ErrUnavailable):
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repoData)
}
