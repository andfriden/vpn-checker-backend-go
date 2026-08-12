package api

import (
	"encoding/json"
	"net/http"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func (h *Handler) ExportJSON(
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

	configs := make(
		[]*model.VPNConfig,
		0,
	)

	for _, result := range results {

		if !result.Success {
			continue
		}

		if result.Config == nil {
			continue
		}

		configs = append(
			configs,
			result.Config,
		)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=working.json",
	)

	json.NewEncoder(w).Encode(
		configs,
	)
}
