package singbox

import (
	"encoding/json"
	"testing"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func TestBuildVLESS(t *testing.T) {
	cfg := &model.VPNConfig{
		Protocol: model.ProtocolVLESS,
		Address:  "example.com",
		Port:     443,
		UUID:     "550e8400-e29b-41d4-a716-446655440000",

		TLS: model.TLSConfig{
			Enabled:    true,
			ServerName: "example.com",
		},

		Transport: model.TransportConfig{
			Type: "ws",
			Path: "/vpn",
			Host: "example.com",
		},
	}

	result, err := Build(cfg, 1080)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(result.Inbounds) != 1 {
		t.Fatalf("inbounds = %d", len(result.Inbounds))
	}

	if result.Inbounds[0].ListenPort != 1080 {
		t.Fatalf("SOCKS port = %d", result.Inbounds[0].ListenPort)
	}

	if len(result.Outbounds) != 2 {
		t.Fatalf("outbounds = %d", len(result.Outbounds))
	}

	out := result.Outbounds[0]

	if out.Type != "vless" {
		t.Fatalf("type = %q", out.Type)
	}

	if out.Server != "example.com" {
		t.Fatalf("server = %q", out.Server)
	}

	if out.ServerPort != 443 {
		t.Fatalf("server port = %d", out.ServerPort)
	}

	if out.UUID != cfg.UUID {
		t.Fatalf("UUID = %q", out.UUID)
	}

	if out.TLS == nil {
		t.Fatal("TLS is nil")
	}

	if out.Transport == nil {
		t.Fatal("transport is nil")
	}

	data, err := Marshal(result)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded map[string]any

	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestBuildHysteria2(t *testing.T) {
	cfg := &model.VPNConfig{
		Protocol: model.ProtocolHysteria2,
		Address:  "example.com",
		Port:     443,
		Password: "secret",

		TLS: model.TLSConfig{
			Enabled:    true,
			ServerName: "example.com",
		},
	}

	result, err := Build(cfg, 1081)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	out := result.Outbounds[0]

	if out.Type != "hysteria2" {
		t.Fatalf("type = %q", out.Type)
	}

	if out.Password != "secret" {
		t.Fatalf("password = %q", out.Password)
	}

	if out.TLS == nil {
		t.Fatal("TLS is nil")
	}
}

func TestBuildTrojan(t *testing.T) {
	cfg := &model.VPNConfig{
		Protocol: model.ProtocolTrojan,
		Address:  "example.com",
		Port:     443,
		Password: "secret",
	}

	result, err := Build(cfg, 1082)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	out := result.Outbounds[0]

	if out.Type != "trojan" {
		t.Fatalf("type = %q", out.Type)
	}

	if out.Password != "secret" {
		t.Fatalf("password = %q", out.Password)
	}

	if out.TLS == nil {
		t.Fatal("TLS is nil")
	}
}

func TestBuildVMess(t *testing.T) {
	cfg := &model.VPNConfig{
		Protocol: model.ProtocolVMess,
		Address:  "example.com",
		Port:     443,
		UUID:     "550e8400-e29b-41d4-a716-446655440000",
	}

	result, err := Build(cfg, 1083)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	out := result.Outbounds[0]

	if out.Type != "vmess" {
		t.Fatalf("type = %q", out.Type)
	}

	if out.UUID != cfg.UUID {
		t.Fatalf("UUID = %q", out.UUID)
	}
}
