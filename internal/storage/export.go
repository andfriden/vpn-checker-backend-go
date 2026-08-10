package storage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func (s *FileStorage) ExportWorking(configs []*model.VPNConfig) error {
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return err
	}

	var lines []string
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		raw := strings.TrimSpace(cfg.Raw)
		if raw != "" {
			lines = append(lines, raw)
		}
	}

	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}

	file := filepath.Join(s.Dir, "all-working.txt")
	return os.WriteFile(file, []byte(content), 0644)
}
