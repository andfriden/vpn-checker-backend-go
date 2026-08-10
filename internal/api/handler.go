package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/app"
	"github.com/andfriden/vpn-checker-backend-go/internal/model"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

type Handler struct {
	runner  *app.Runner
	storage *storage.FileStorage
}

type APIResult struct {
	Protocol string `json:"protocol,omitempty"`
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port,omitempty"`
	Latency  int64  `json:"latency_ms"`
	IP       string `json:"ip,omitempty"`
	Success  bool   `json:"success"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type ResultsResponse struct {
	Data       []APIResult `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func NewHandler(
	runner *app.Runner,
	storage *storage.FileStorage,
) *Handler {
	return &Handler{
		runner:  runner,
		storage: storage,
	}
}

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

func (h *Handler) Results(
	w http.ResponseWriter,
	r *http.Request,
) {
	results, err := h.storage.LoadResults()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	query := r.URL.Query()

	workingFilter := strings.ToLower(
		strings.TrimSpace(
			query.Get("working"),
		),
	)

	protocolFilter := strings.ToLower(
		strings.TrimSpace(
			query.Get("protocol"),
		),
	)

	sortFilter := strings.ToLower(
		strings.TrimSpace(
			query.Get("sort"),
		),
	)

	filtered := filterResults(
		results,
		workingFilter,
		protocolFilter,
	)

	sortResults(
		filtered,
		sortFilter,
	)

	total := len(filtered)

	page := parsePositiveInt(
		query.Get("page"),
		1,
	)

	limit := parsePositiveInt(
		query.Get("limit"),
		20,
	)

	pages := 0

	if total > 0 {
		pages = (total + limit - 1) / limit
	}

	if pages > 0 && page > pages {
		page = pages
	}

	start := (page - 1) * limit

	if start > total {
		start = total
	}

	end := start + limit

	if end > total {
		end = total
	}

	pageResults := filtered[start:end]

	response := ResultsResponse{
		Data: toAPIResults(pageResults),
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: pages,
		},
	}

	writeJSON(w, response)
}

func (h *Handler) Best(
	w http.ResponseWriter,
	r *http.Request,
) {
	results, err := h.storage.LoadResults()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	query := r.URL.Query()

	protocolFilter := strings.ToLower(
		strings.TrimSpace(
			query.Get("protocol"),
		),
	)

	limit := parsePositiveInt(
		query.Get("limit"),
		10,
	)

	filtered := filterResults(
		results,
		"true",
		protocolFilter,
	)

	sortResults(
		filtered,
		"latency",
	)

	filtered = deduplicateResults(filtered)

	if limit < len(filtered) {
		filtered = filtered[:limit]
	}

	writeJSON(
		w,
		toAPIResults(filtered),
	)
}

func (h *Handler) Stats(
	w http.ResponseWriter,
	r *http.Request,
) {
	stats, err := h.storage.LoadStats()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, stats)
}

func filterResults(
	results []model.CheckResult,
	workingFilter string,
	protocolFilter string,
) []model.CheckResult {
	filtered := make(
		[]model.CheckResult,
		0,
		len(results),
	)

	for _, result := range results {
		switch workingFilter {
		case "true":
			if !result.Success {
				continue
			}

		case "false":
			if result.Success {
				continue
			}
		}

		if protocolFilter != "" {
			if result.Config == nil {
				continue
			}

			protocol := strings.ToLower(
				string(result.Config.Protocol),
			)

			if protocol != protocolFilter {
				continue
			}
		}

		filtered = append(
			filtered,
			result,
		)
	}

	return filtered
}

func sortResults(
	results []model.CheckResult,
	sortFilter string,
) {
	switch sortFilter {
	case "latency":
		sort.SliceStable(
			results,
			func(i, j int) bool {
				return results[i].Latency <
					results[j].Latency
			},
		)

	case "latency_desc":
		sort.SliceStable(
			results,
			func(i, j int) bool {
				return results[i].Latency >
					results[j].Latency
			},
		)
	}
}

func deduplicateResults(
	results []model.CheckResult,
) []model.CheckResult {
	seen := make(
		map[string]struct{},
		len(results),
	)

	unique := make(
		[]model.CheckResult,
		0,
		len(results),
	)

	for _, result := range results {
		if result.Config == nil {
			continue
		}

		key := strings.ToLower(
			string(result.Config.Protocol),
		) + "|" +
			strings.ToLower(
				result.Config.Address,
			) + "|" +
			strconv.Itoa(
				result.Config.Port,
			)

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, result)
	}

	return unique
}

func toAPIResults(
	results []model.CheckResult,
) []APIResult {
	data := make(
		[]APIResult,
		0,
		len(results),
	)

	for _, result := range results {
		apiResult := APIResult{
			Latency: result.Latency.Milliseconds(),
			IP:      result.IP,
			Success: result.Success,
		}

		if result.Config != nil {
			apiResult.Protocol = string(
				result.Config.Protocol,
			)

			apiResult.Address = result.Config.Address
			apiResult.Port = result.Config.Port
		}

		data = append(
			data,
			apiResult,
		)
	}

	return data
}

func parsePositiveInt(
	value string,
	fallback int,
) int {
	value = strings.TrimSpace(value)

	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)

	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func writeJSON(
	w http.ResponseWriter,
	value interface{},
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) CheckStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		h.runner.Status(),
	)
}
