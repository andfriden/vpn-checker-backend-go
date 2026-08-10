package downloader

import (
	"strings"
)

type Downloader struct {
}

func New() *Downloader {

	return &Downloader{}
}

func (d *Downloader) ParseList(
	data string,
) []string {

	lines := strings.Split(
		data,
		"\n",
	)

	seen := make(map[string]bool)

	var result []string

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if seen[line] {
			continue
		}

		seen[line] = true

		result = append(
			result,
			line,
		)
	}

	return result
}
