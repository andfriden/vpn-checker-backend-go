package downloader

import (
	"strings"
)

func (d *Downloader) ParseList(data string) []string {

	lines := strings.Split(data, "\n")

	result := make([]string, 0)

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		// убираем комментарии
		if strings.HasPrefix(line, "#") {
			continue
		}

		result = append(result, line)
	}

	return result
}
