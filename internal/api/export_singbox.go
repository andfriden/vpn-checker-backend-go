package api

import (
	"encoding/json"
	"net/http"

	"github.com/andfriden/vpn-checker-backend-go/internal/singbox"
)

func (h *Handler) ExportSingBox(
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

	config := singbox.Config{
		Log: singbox.LogConfig{
			Level: "error",
		},

		Outbounds: []singbox.OutboundConfig{
			{
				Type: "direct",
				Tag:  "direct",
			},
		},

		Route: singbox.RouteConfig{
			Final: "direct",
		},
	}

	for _, result := range results {

		if !result.Success {
			continue
		}

		if result.Config == nil {
			continue
		}

		sb, err := singbox.Build(
			result.Config,
			1080,
		)

		if err != nil {
			continue
		}

		for _, outbound := range sb.Outbounds {

			if outbound.Tag == singbox.VPNOutboundTag {

				config.Outbounds = append(
					config.Outbounds,
					outbound,
				)

			}
		}
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=singbox.json",
	)

	json.NewEncoder(w).Encode(
		config,
	)
}
