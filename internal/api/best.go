package api

import (
	"net/http"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/service"
)

func (h *Handler) Best(
	w http.ResponseWriter,
	r *http.Request,
) {

	query := r.URL.Query()

	protocol := strings.ToLower(
		strings.TrimSpace(
			query.Get("protocol"),
		),
	)

	limit := service.ParseInt(
		query.Get("limit"),
		10,
	)

	results, err := h.resultsService.Best(
		protocol,
		limit,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	writeJSON(
		w,
		convertAPIResults(
			results,
		),
	)
}
