package api

import "net/http"

func (h *Handler) Stats(
	w http.ResponseWriter,
	r *http.Request,
) {

	stats, err := h.resultsService.Stats()

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
		stats,
	)
}
