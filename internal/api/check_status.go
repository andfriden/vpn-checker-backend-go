package api

import (
	"net/http"
)

func (h *Handler) CheckStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		h.runner.Status(),
	)
}
