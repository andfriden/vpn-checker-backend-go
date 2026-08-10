package api

import "net/http"

func Routes(h *Handler) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		h.Health,
	)

	mux.HandleFunc(
		"/api/results",
		h.Results,
	)

	mux.HandleFunc(
		"/api/check",
		h.Check,
	)

	mux.HandleFunc(
		"/api/stats",
		h.Stats,
	)

	return mux
}
