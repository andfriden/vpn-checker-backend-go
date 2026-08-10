package singbox

import (
	"encoding/json"
	"fmt"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

const (
	SOCKSInboundTag = "socks-in"
	VPNOutboundTag  = "vpn-out"
	DirectTag       = "direct"
	DNSDirectTag    = "dns-direct"
	DNSRemoteTag    = "dns-remote"
)

func Build(cfg *model.VPNConfig, socksPort int) (*Config, error) {
	if cfg == nil {
		return nil, fmt.Errorf("VPN config is nil")
	}

	if cfg.Address == "" {
		return nil, fmt.Errorf("VPN address is empty")
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid VPN port: %d", cfg.Port)
	}

	if socksPort < 1 || socksPort > 65535 {
		return nil, fmt.Errorf("invalid SOCKS port: %d", socksPort)
	}

	outbound, err := buildOutbound(cfg)
	if err != nil {
		return nil, err
	}

	return &Config{
		Log: LogConfig{
			Level: "error",
		},

		DNS: &DNSConfig{
			Servers: []DNSServer{
				{
					Type: "local",
					Tag:  DNSDirectTag,
				},
				{
					Type:   "https",
					Tag:    DNSRemoteTag,
					Server: "8.8.8.8",
					Detour: VPNOutboundTag,
				},
			},
		},

		Inbounds: []InboundConfig{
			{
				Type:       "mixed",
				Tag:        SOCKSInboundTag,
				Listen:     "127.0.0.1",
				ListenPort: socksPort,
			},
		},

		Outbounds: []OutboundConfig{
			*outbound,
			{
				Type: "direct",
				Tag:  DirectTag,
			},
		},

		Route: RouteConfig{
			Final:                 VPNOutboundTag,
			DefaultDomainResolver: DNSDirectTag,
		},
	}, nil
}

func buildOutbound(cfg *model.VPNConfig) (*OutboundConfig, error) {
	out := &OutboundConfig{
		Tag:        VPNOutboundTag,
		Server:     cfg.Address,
		ServerPort: cfg.Port,
	}

	switch cfg.Protocol {

	case model.ProtocolHysteria2:
		out.Type = "hysteria2"

		if cfg.Password == "" {
			return nil, fmt.Errorf("Hysteria2 password is empty")
		}

		out.Password = cfg.Password

		if cfg.UpMbps > 0 {
			out.UpMbps = cfg.UpMbps
		}

		if cfg.DownMbps > 0 {
			out.DownMbps = cfg.DownMbps
		}

		/*
			Hysteria2 в sing-box использует собственные
			поля outbound и НЕ использует transport.

			Поэтому здесь намеренно нет:

			    out.Transport = ...

			Иначе sing-box 1.12 выдаёт:

			    unknown field "transport"
		*/

		if cfg.TLS.Enabled {
			out.TLS = buildTLS(cfg)
		}

	case model.ProtocolVLESS:
		out.Type = "vless"

		if cfg.UUID == "" {
			return nil, fmt.Errorf("VLESS UUID is empty")
		}

		out.UUID = cfg.UUID

		if cfg.TLS.Enabled {
			out.TLS = buildTLS(cfg)
		}

		out.Transport = buildTransport(cfg)

	case model.ProtocolVMess:
		out.Type = "vmess"

		if cfg.UUID == "" {
			return nil, fmt.Errorf("VMess UUID is empty")
		}

		out.UUID = cfg.UUID

		if cfg.TLS.Enabled {
			out.TLS = buildTLS(cfg)
		}

		out.Transport = buildTransport(cfg)

	case model.ProtocolTrojan:
		out.Type = "trojan"

		if cfg.Password == "" {
			return nil, fmt.Errorf("Trojan password is empty")
		}

		out.Password = cfg.Password

		/*
			Trojan всегда работает поверх TLS
			в нашей модели.
		*/
		out.TLS = buildTLS(cfg)

		out.Transport = buildTransport(cfg)

	default:
		return nil, fmt.Errorf(
			"unsupported protocol: %s",
			cfg.Protocol,
		)
	}

	return out, nil
}

func buildTLS(cfg *model.VPNConfig) *TLSConfig {
	tls := &TLSConfig{
		Enabled:    true,
		ServerName: cfg.TLS.ServerName,
		Insecure:   cfg.TLS.Insecure,
	}

	if tls.ServerName == "" {
		tls.ServerName = cfg.Address
	}

	if cfg.TLS.Fingerprint != "" {
		tls.UTLS = &UTLSConfig{
			Enabled:     true,
			Fingerprint: cfg.TLS.Fingerprint,
		}
	}

	if cfg.TLS.PublicKey != "" || cfg.TLS.ShortID != "" {
		tls.Reality = &RealityConfig{
			Enabled:   true,
			PublicKey: cfg.TLS.PublicKey,
			ShortID:   cfg.TLS.ShortID,
		}
	}

	return tls
}

func buildTransport(cfg *model.VPNConfig) *TransportConfig {
	switch cfg.Transport.Type {

	case "":
		return nil

	case "ws":
		transport := &TransportConfig{
			Type: cfg.Transport.Type,
			Path: cfg.Transport.Path,
		}

		if cfg.Transport.Host != "" {
			transport.Headers = map[string]string{
				"Host": cfg.Transport.Host,
			}
		}

		return transport

	case "grpc":
		return &TransportConfig{
			Type:        "grpc",
			ServiceName: cfg.Transport.Service,
		}

	default:
		return nil
	}
}

func Marshal(cfg *Config) ([]byte, error) {
	if cfg == nil {
		return nil, fmt.Errorf("sing-box config is nil")
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf(
			"marshal sing-box config: %w",
			err,
		)
	}

	return data, nil
}
