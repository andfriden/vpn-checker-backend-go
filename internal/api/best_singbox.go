package api

import (
	"encoding/json"
	"net/http"

	"github.com/andfriden/vpn-checker-backend-go/internal/service"
	"github.com/andfriden/vpn-checker-backend-go/internal/singbox"
)

func (h *Handler) BestSingBox(
	w http.ResponseWriter,
	r *http.Request,
) {
	results, err := h.resultsService.All()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	working := service.Filter(
		results,
		"true",
		"",
	)

	service.Sort(
		working,
		"latency",
	)

	if len(working) == 0 {
		http.Error(
			w,
			"no working configs",
			http.StatusNotFound,
		)
		return
	}

	best := working[0]

	config, err := singbox.Build(
		best.Config,
		1080,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=best-singbox.json",
	)

	_ = json.NewEncoder(w).Encode(config)
}
