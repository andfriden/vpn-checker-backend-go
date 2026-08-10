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

	url := "https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/main/BLACK_SS%2BAll_RUS.txt"

	started := h.runner.RunAsync(
		r.Context(),
		url,
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
		},
	)
}
