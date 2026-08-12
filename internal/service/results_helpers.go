package service

import (
	"sort"
	"strconv"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func Filter(
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

func Sort(
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

func Deduplicate(
	results []model.CheckResult,
) []model.CheckResult {

	seen := make(
		map[string]struct{},
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
		) +
			"|" +
			strings.ToLower(
				result.Config.Address,
			) +
			"|" +
			strconv.Itoa(
				result.Config.Port,
			)

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}

		unique = append(
			unique,
			result,
		)
	}

	return unique
}

func convertResults(
	results []model.CheckResult,
) []ResultItem {

	data := make(
		[]ResultItem,
		0,
		len(results),
	)

	for _, result := range results {

		item := ResultItem{
			Latency: result.Latency.Milliseconds(),
			IP:      result.IP,
			Success: result.Success,
		}

		if result.Config != nil {

			item.Protocol =
				string(result.Config.Protocol)

			item.Address =
				result.Config.Address

			item.Port =
				result.Config.Port
		}

		data = append(
			data,
			item,
		)
	}

	return data
}

func ParseInt(
	value string,
	fallback int,
) int {

	value = strings.TrimSpace(
		value,
	)

	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(
		value,
	)

	if err != nil || result <= 0 {
		return fallback
	}

	return result
}
