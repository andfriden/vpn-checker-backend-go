package api

import (
	"context"
	"encoding/json"
	"net/http"
)

func (h *Handler) Check(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}

	source := r.URL.Query().Get("source")

	var url string

	switch source {
	case "white":
		url = h.whiteURL

	case "black":
		url = h.blackURL

	default:
		source = "black"
		url = h.blackURL
	}

	if url == "" {
		http.Error(
			w,
			"source url is not configured",
			http.StatusBadRequest,
		)

		return
	}

	started := h.runner.RunAsync(
		context.Background(),
		url,
		source,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if !started {
		json.NewEncoder(w).Encode(
			map[string]string{
				"status": "already_running",
			},
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "started",
			"source": source,
		},
	)
}
