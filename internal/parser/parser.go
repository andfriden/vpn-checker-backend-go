package parser

import (
	"fmt"
	"strings"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

type Parser interface {
	Parse(raw string) (*model.VPNConfig, error)
	Protocol() model.Protocol
}

func Parse(raw string) (*model.VPNConfig, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return nil, fmt.Errorf("empty VPN config")
	}

	switch {
	case strings.HasPrefix(strings.ToLower(raw), "vless://"):
		return parseVLESS(raw)

	case strings.HasPrefix(strings.ToLower(raw), "vmess://"):
		return parseVMess(raw)

	case strings.HasPrefix(strings.ToLower(raw), "trojan://"):
		return parseTrojan(raw)

	case strings.HasPrefix(strings.ToLower(raw), "hysteria2://"):
		return parseHysteria2(raw)

	case strings.HasPrefix(strings.ToLower(raw), "hy2://"):
		return parseHysteria2(raw)

	default:
		return nil, fmt.Errorf("unsupported protocol")
	}
}
