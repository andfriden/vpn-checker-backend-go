package extractor

import (
	"strings"
)

var protocols = []string{
	"vless://",
	"vmess://",
	"trojan://",
	"ss://",
	"ssconf://",
	"hysteria://",
	"hysteria2://",
	"hy2://",
	"tuic://",
	"juicity://",
	"naive+https://",
}

func Extract(data string) []string {
	lines := strings.Split(data, "\n")

	result := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		for _, protocol := range protocols {
			if strings.HasPrefix(strings.ToLower(line), protocol) {
				result = append(result, line)
				break
			}
		}
	}

	return result
}
