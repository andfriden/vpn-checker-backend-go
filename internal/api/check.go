package api

import (
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


	url := h.blackURL


	if source == "white" {
		url = h.whiteURL
	}


	started := h.runner.RunAsync(
		r.Context(),
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
