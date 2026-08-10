package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func parseTrojan(raw string) (*model.VPNConfig, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse Trojan URI: %w", err)
	}

	if strings.ToLower(u.Scheme) != "trojan" {
		return nil, fmt.Errorf("invalid Trojan scheme: %q", u.Scheme)
	}

	if u.Hostname() == "" {
		return nil, fmt.Errorf("Trojan address is empty")
	}

	if u.User == nil {
		return nil, fmt.Errorf("Trojan password is missing")
	}

	password := u.User.Username()

	if password == "" {
		return nil, fmt.Errorf("Trojan password is empty")
	}

	port := 443

	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid Trojan port: %w", err)
		}
	}

	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid Trojan port: %d", port)
	}

	query := u.Query()

	config := &model.VPNConfig{
		Raw:      raw,
		Protocol: model.ProtocolTrojan,
		Address:  u.Hostname(),
		Port:     port,
		Password: password,
		Name:     fragmentName(u.Fragment),
		TLS: model.TLSConfig{
			Enabled:     true,
			ServerName:  query.Get("sni"),
			Fingerprint: query.Get("fp"),
		},
	}

	if config.TLS.ServerName == "" {
		config.TLS.ServerName = u.Hostname()
	}

	switch strings.ToLower(query.Get("type")) {
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
