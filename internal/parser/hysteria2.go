package parser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func parseHysteria2(raw string) (*model.VPNConfig, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, err
	}

	password := ""
	if u.User != nil {
		password = u.User.Username()

		if pass, ok := u.User.Password(); ok && pass != "" {
			password += ":" + pass
		}
	}

	cfg := &model.VPNConfig{
		Raw:      raw,
		Protocol: model.ProtocolHysteria2,
		Address:  u.Hostname(),
		Port:     port,
		Password: password,
		Name:     fragmentName(u.Fragment),
	}
	// TLS
	cfg.TLS.Enabled = true

	query := u.Query()

	if sni := query.Get("sni"); sni != "" {
		cfg.TLS.ServerName = sni
	}

	if query.Get("insecure") == "1" ||
		strings.EqualFold(query.Get("allowInsecure"), "true") {
		cfg.TLS.Insecure = true
	}

	// Hysteria2 transport
	// В текущей модели используется Transport
	cfg.Transport.Type = "quic"

	return cfg, nil
}
