package source

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sources []string `yaml:"sources"`
}

func Load(path string) ([]string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(
		data,
		&cfg,
	)

	if err != nil {
		return nil, err
	}

	return cfg.Sources, nil
}
