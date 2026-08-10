package api

import (
	"encoding/json"
	"net/http"

	"github.com/andfriden/vpn-checker-backend-go/internal/app"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

type Handler struct {
	runner  *app.Runner
	storage *storage.FileStorage
}

func NewHandler(
	runner *app.Runner,
	storage *storage.FileStorage,
) *Handler {

	return &Handler{
		runner:  runner,
		storage: storage,
	}
}

func (h *Handler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "ok",
		},
	)
}

func (h *Handler) Results(
	w http.ResponseWriter,
	r *http.Request,
) {

	results, err := h.storage.LoadResults()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		results,
	)
}

func (h *Handler) Stats(
	w http.ResponseWriter,
	r *http.Request,
) {

	stats, err := h.storage.LoadStats()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		stats,
	)
}
