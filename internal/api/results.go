package api

import (
	"net/http"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/service"
)

func (h *Handler) Results(
	w http.ResponseWriter,
	r *http.Request,
) {
	query := r.URL.Query()

	result, err := h.resultsService.List(
		service.ResultsQuery{
			Working:  strings.ToLower(strings.TrimSpace(query.Get("working"))),
			Protocol: strings.ToLower(strings.TrimSpace(query.Get("protocol"))),
			Sort:     strings.ToLower(strings.TrimSpace(query.Get("sort"))),
			Page:     service.ParseInt(query.Get("page"), 1),
			Limit:    service.ParseInt(query.Get("limit"), 20),
		},
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, result)
}
