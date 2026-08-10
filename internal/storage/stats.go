package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func (s *FileStorage) SaveStats(
	stats model.Stats,
) error {

	if err := os.MkdirAll(
		s.Dir,
		0755,
	); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		stats,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		filepath.Join(
			s.Dir,
			"stats.json",
		),
		data,
		0644,
	)
}

func (s *FileStorage) LoadStats() (
	model.Stats,
	error,
) {

	var stats model.Stats

	data, err := os.ReadFile(
		filepath.Join(
			s.Dir,
			"stats.json",
		),
	)

	if err != nil {
		return stats, err
	}

	err = json.Unmarshal(
		data,
		&stats,
	)

	return stats, err
}
