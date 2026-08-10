package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Checker CheckerConfig `yaml:"checker"`
	SingBox SingBoxConfig `yaml:"singbox"`

	Storage StorageConfig `yaml:"storage"`

	Cleanup CleanupConfig `yaml:"cleanup"`
	Metrics MetricsConfig `yaml:"metrics"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type CheckerConfig struct {
	Workers             int           `yaml:"workers"`
	MaxConcurrentChecks int           `yaml:"max_concurrent_checks"`
	Timeout             time.Duration `yaml:"timeout"`
	IPCheckURL          string        `yaml:"ip_check_url"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

type SingBoxConfig struct {
	Binary          string        `yaml:"binary"`
	StartupTimeout  time.Duration `yaml:"startup_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	SocksHost       string        `yaml:"socks_host"`
}

type StorageConfig struct {
	Path string `yaml:"path"`
}

type CleanupConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Interval  time.Duration `yaml:"interval"`
	DeadAfter int           `yaml:"dead_after"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	if c.Checker.Workers < 1 {
		return fmt.Errorf("checker.workers must be greater than 0")
	}

	if c.Checker.MaxConcurrentChecks < 1 {
		return fmt.Errorf("checker.max_concurrent_checks must be greater than 0")
	}

	if c.Checker.Timeout <= 0 {
		return fmt.Errorf("checker.timeout must be greater than 0")
	}

	if c.Checker.IPCheckURL == "" {
		return fmt.Errorf("checker.ip_check_url cannot be empty")
	}

	if c.SingBox.Binary == "" {
		return fmt.Errorf("singbox.binary cannot be empty")
	}

	if c.Storage.Path == "" {
		return fmt.Errorf("storage.path cannot be empty")
	}

	return nil
}
