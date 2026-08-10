package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

type FileStorage struct {
	Dir string
}

func New(dir string) *FileStorage {

	return &FileStorage{
		Dir: dir,
	}
}

func (s *FileStorage) SaveResults(results []model.CheckResult) error {

	err := os.MkdirAll(s.Dir, 0755)

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		results,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	file := filepath.Join(
		s.Dir,
		"results.json",
	)

	return os.WriteFile(
		file,
		data,
		0644,
	)
}

func (s *FileStorage) LoadResults() ([]model.CheckResult, error) {

	file := filepath.Join(
		s.Dir,
		"results.json",
	)

	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var results []model.CheckResult

	err = json.Unmarshal(
		data,
		&results,
	)

	if err != nil {
		return nil, err
	}

	return results, nil
}
