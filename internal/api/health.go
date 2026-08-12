package api

import "net/http"

func (h *Handler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		map[string]string{
			"status": "ok",
		},
	)
}
