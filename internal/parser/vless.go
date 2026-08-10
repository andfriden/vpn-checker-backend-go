package parser

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func parseVLESS(raw string) (*model.VPNConfig, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse VLESS URI: %w", err)
	}

	if u.Scheme != "vless" {
		return nil, fmt.Errorf("invalid VLESS scheme: %q", u.Scheme)
	}

	if u.Hostname() == "" {
		return nil, fmt.Errorf("VLESS address is empty")
	}

	port := 443

	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid VLESS port: %w", err)
		}
	}

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid VLESS port: %d", port)
	}

	if u.User == nil {
		return nil, fmt.Errorf("VLESS UUID is missing")
	}

	uuid := u.User.Username()

	if uuid == "" {
		return nil, fmt.Errorf("VLESS UUID is empty")
	}

	config := &model.VPNConfig{
		Raw:      raw,
		Protocol: model.ProtocolVLESS,
		Address:  u.Hostname(),
		Port:     port,
		UUID:     uuid,
		Name:     fragmentName(u.Fragment),
		TLS: model.TLSConfig{
			Enabled: false,
		},
	}

	query := u.Query()

	security := query.Get("security")

	switch security {
	case "tls":
		config.TLS.Enabled = true

	case "reality":
		config.TLS.Enabled = true
		config.TLS.Fingerprint = query.Get("fp")
		config.TLS.PublicKey = query.Get("pbk")
		config.TLS.ShortID = query.Get("sid")
	}

	if serverName := query.Get("sni"); serverName != "" {
		config.TLS.ServerName = serverName
	}

	switch query.Get("type") {
	case "ws":
		config.Transport.Type = "ws"
		config.Transport.Path = query.Get("path")
		config.Transport.Host = query.Get("host")

	case "grpc":
		config.Transport.Type = "grpc"
		config.Transport.Service = query.Get("serviceName")

	case "tcp", "":
		config.Transport.Type = "tcp"

	default:
		config.Transport.Type = query.Get("type")
	}

	return config, nil
}
