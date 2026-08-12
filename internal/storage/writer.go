package storage

import (
	"os"
	"path/filepath"
	"strings"
)

func Save(path string, configs []string) error {

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := strings.Join(
		configs,
		"\n",
	)

	return os.WriteFile(
		path,
		[]byte(data),
		0644,
	)
}
