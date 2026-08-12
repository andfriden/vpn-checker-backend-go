package pipeline

import (
	"sort"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/decoder"
	"github.com/andfriden/vpn-checker-backend-go/internal/downloader"
	"github.com/andfriden/vpn-checker-backend-go/internal/extractor"
	"github.com/andfriden/vpn-checker-backend-go/internal/normalizer"
	"github.com/andfriden/vpn-checker-backend-go/internal/parser"
	"github.com/andfriden/vpn-checker-backend-go/internal/validator"
)

func Run(urls []string) ([]string, error) {

	unique := make(map[string]struct{})

	for _, url := range urls {

		body, err := downloader.DownloadURL(url)
		if err != nil {
			continue
		}

		for _, item := range body {

			decoded := decoder.Decode(item)

			configs := extractor.Extract(decoded)

			for _, config := range configs {

				config = strings.TrimSpace(config)

				if config == "" {
					continue
				}

				if !validator.Validate(config) {
					continue
				}

				_, err := parser.Parse(config)

				if err != nil {
					continue
				}

				unique[config] = struct{}{}
			}
		}
	}

	configs := make([]string, 0, len(unique))

	for config := range unique {
		configs = append(configs, config)
	}

	// смысловая очистка дублей
	configs = normalizer.Unique(configs)

	sort.Strings(configs)

	return configs, nil
}
