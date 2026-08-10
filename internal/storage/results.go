package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/checker"
)

type Storage struct {
	Path string
}

func (s *Storage) Init() error {

	return os.MkdirAll(
		s.Path,
		0755,
	)
}

func (s *Storage) SaveWorking(
	results []*checker.Result,
) error {

	if err := s.Init(); err != nil {
		return err
	}

	var working []*checker.Result

	for _, r := range results {

		if r == nil {
			continue
		}

		if r.Success {

			working = append(
				working,
				r,
			)

		}
	}

	data, err := json.MarshalIndent(
		working,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	file := filepath.Join(
		s.Path,
		"working.json",
	)

	return os.WriteFile(
		file,
		data,
		0644,
	)
}

func (s *Storage) SaveFailed(
	results []*checker.Result,
) error {

	if err := s.Init(); err != nil {
		return err
	}

	var failed []*checker.Result

	for _, r := range results {

		if r == nil {
			continue
		}

		if !r.Success {

			failed = append(
				failed,
				r,
			)
		}
	}

	data, err := json.MarshalIndent(
		failed,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	file := filepath.Join(
		s.Path,
		"failed.json",
	)

	return os.WriteFile(
		file,
		data,
		0644,
	)
}

func (s *Storage) SaveStats(
	results []*checker.Result,
) error {

	if err := s.Init(); err != nil {
		return err
	}

	total := len(results)

	working := 0

	for _, r := range results {

		if r != nil && r.Success {
			working++
		}
	}

	stats := map[string]interface{}{

		"total": total,

		"working": working,

		"failed": total - working,

		"updated": time.Now(),
	}

	data, err := json.MarshalIndent(
		stats,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	file := filepath.Join(
		s.Path,
		"stats.json",
	)

	return os.WriteFile(
		file,
		data,
		0644,
	)
}

func (s *Storage) SaveAll(
	results []*checker.Result,
) error {

	if err := s.SaveWorking(results); err != nil {
		return fmt.Errorf(
			"save working: %w",
			err,
		)
	}

	if err := s.SaveFailed(results); err != nil {
		return fmt.Errorf(
			"save failed: %w",
			err,
		)
	}

	return s.SaveStats(results)
}
