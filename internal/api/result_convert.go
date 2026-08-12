package api

import (
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

type APIResult struct {
	Protocol string `json:"protocol,omitempty"`
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port,omitempty"`
	Latency  int64  `json:"latency_ms"`
	IP       string `json:"ip,omitempty"`
	Success  bool   `json:"success"`
	Source   string `json:"source,omitempty"`
}

func convertAPIResults(
	results []model.CheckResult,
) []APIResult {
	data := make(
		[]APIResult,
		0,
		len(results),
	)

	for _, result := range results {
		item := APIResult{
			Latency: result.Latency.Milliseconds(),
			IP:      result.IP,
			Success: result.Success,
		}

		if result.Config != nil {
			item.Protocol = string(result.Config.Protocol)
			item.Address = result.Config.Address
			item.Port = result.Config.Port
			item.Source = detectSource(result.Config.Name)
		}

		data = append(data, item)
	}

	return data
}

func detectSource(name string) string {
	name = strings.ToUpper(name)

	switch {
	case strings.Contains(name, "[BL]"):
		return "BLACK"

	case strings.Contains(name, "[WL]"):
		return "WHITE"

	default:
		return ""
	}
}
