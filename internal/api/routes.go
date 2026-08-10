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
		"/api/best",
		h.Best,
	)

	mux.HandleFunc(
		"/api/check",
		h.Check,
	)

	mux.HandleFunc(
		"/api/check/status",
		h.CheckStatus,
	)

	mux.HandleFunc(
		"/api/stats",
		h.Stats,
	)

	mux.HandleFunc(
		"/api/export",
		h.Export,
	)

	return mux
}
