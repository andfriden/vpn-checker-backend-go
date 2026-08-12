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
		"/api/best/singbox",
		h.BestSingBox,
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

	mux.HandleFunc(
		"/api/export/json",
		h.ExportJSON,
	)

	mux.HandleFunc(
		"/api/export/singbox",
		h.ExportSingBox,
	)

	return mux
}
