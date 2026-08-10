package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

type vmessConfig struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port any    `json:"port"`
	ID   string `json:"id"`
	Aid  any    `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
	SNI  string `json:"sni"`
	FP   string `json:"fp"`
}

func parseVMess(raw string) (*model.VPNConfig, error) {
	const prefix = "vmess://"

	encoded := strings.TrimSpace(raw)

	if !strings.HasPrefix(strings.ToLower(encoded), prefix) {
		return nil, fmt.Errorf("invalid VMess URI")
	}

	encoded = encoded[len(prefix):]

	if encoded == "" {
		return nil, fmt.Errorf("VMess payload is empty")
	}

	data, err := decodeVMessPayload(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode VMess payload: %w", err)
	}

	var vm vmessConfig

	if err := json.Unmarshal(data, &vm); err != nil {
		return nil, fmt.Errorf("parse VMess JSON: %w", err)
	}

	if vm.Add == "" {
		return nil, fmt.Errorf("VMess address is empty")
	}

	if vm.ID == "" {
		return nil, fmt.Errorf("VMess UUID is empty")
	}

	port, err := parseVMessPort(vm.Port)
	if err != nil {
		return nil, err
	}

	config := &model.VPNConfig{
		Raw:      raw,
		Protocol: model.ProtocolVMess,
		Address:  vm.Add,
		Port:     port,
		UUID:     vm.ID,
		Name:     vm.Ps,
		TLS: model.TLSConfig{
			Enabled:     false,
			ServerName:  vm.SNI,
			Fingerprint: vm.FP,
		},
		Transport: model.TransportConfig{
			Type: vm.Net,
			Host: vm.Host,
			Path: vm.Path,
		},
	}

	if vm.TLS != "" && strings.ToLower(vm.TLS) != "none" {
		config.TLS.Enabled = true
	}

	if config.Transport.Type == "" {
		config.Transport.Type = "tcp"
	}

	switch strings.ToLower(config.Transport.Type) {
	case "ws":
		config.Transport.Type = "ws"

	case "grpc":
		config.Transport.Type = "grpc"

	case "http":
		config.Transport.Type = "http"

	case "h2":
		config.Transport.Type = "http"

	case "tcp":
		config.Transport.Type = "tcp"
	}

	return config, nil
}

func decodeVMessPayload(payload string) ([]byte, error) {
	payload = strings.TrimSpace(payload)

	decoders := []func(string) ([]byte, error){
		func(s string) ([]byte, error) {
			return base64.StdEncoding.DecodeString(s)
		},
		func(s string) ([]byte, error) {
			return base64.RawStdEncoding.DecodeString(s)
		},
		func(s string) ([]byte, error) {
			return base64.URLEncoding.DecodeString(s)
		},
		func(s string) ([]byte, error) {
			return base64.RawURLEncoding.DecodeString(s)
		},
	}

	var lastErr error

	for _, decode := range decoders {
		data, err := decode(payload)
		if err == nil {
			return data, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

func parseVMessPort(value any) (int, error) {
	switch v := value.(type) {
	case float64:
		port := int(v)

		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("invalid VMess port: %d", port)
		}

		return port, nil

	case string:
		if v == "" {
			return 0, fmt.Errorf("VMess port is empty")
		}

		port, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid VMess port: %w", err)
		}

		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("invalid VMess port: %d", port)
		}

		return port, nil

	case json.Number:
		port, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("invalid VMess port: %w", err)
		}

		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("invalid VMess port: %d", port)
		}

		return int(port), nil

	default:
		return 0, fmt.Errorf("unsupported VMess port type: %T", value)
	}
}

func vmessNameFromRaw(raw string) string {
	if idx := strings.Index(raw, "#"); idx >= 0 && idx+1 < len(raw) {
		name, err := url.QueryUnescape(raw[idx+1:])
		if err == nil {
			return name
		}
	}

	return ""
}
