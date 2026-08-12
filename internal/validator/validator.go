package validator

import (
	"net/url"
	"strings"
)

var allowedSchemes = map[string]bool{
	"vless":     true,
	"vmess":     true,
	"trojan":    true,
	"ss":        true,
	"hysteria":  true,
	"hysteria2": true,
	"hy2":       true,
}

func Validate(config string) bool {

	config = strings.TrimSpace(config)

	if config == "" {
		return false
	}

	if !strings.Contains(config, "://") {
		return false
	}

	u, err := url.Parse(config)

	if err != nil {
		return false
	}

	if !allowedSchemes[u.Scheme] {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
