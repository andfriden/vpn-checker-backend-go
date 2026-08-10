package parser

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func TestParseVLESS(t *testing.T) {
	raw := "vless://550e8400-e29b-41d4-a716-446655440000@example.com:443?security=tls&sni=example.com&type=ws&path=%2Fvpn&host=example.com#VLESS-Test"

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Protocol != model.ProtocolVLESS {
		t.Fatalf("protocol = %q", cfg.Protocol)
	}

	if cfg.Address != "example.com" {
		t.Fatalf("address = %q", cfg.Address)
	}

	if cfg.Port != 443 {
		t.Fatalf("port = %d", cfg.Port)
	}

	if cfg.UUID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("UUID = %q", cfg.UUID)
	}

	if !cfg.TLS.Enabled {
		t.Fatal("TLS should be enabled")
	}

	if cfg.TLS.ServerName != "example.com" {
		t.Fatalf("SNI = %q", cfg.TLS.ServerName)
	}

	if cfg.Transport.Type != "ws" {
		t.Fatalf("transport = %q", cfg.Transport.Type)
	}

	if cfg.Transport.Path != "/vpn" {
		t.Fatalf("path = %q", cfg.Transport.Path)
	}

	if cfg.Name != "VLESS-Test" {
		t.Fatalf("name = %q", cfg.Name)
	}
}

func TestParseVMess(t *testing.T) {
	payload := map[string]any{
		"v":    "2",
		"ps":   "VMess-Test",
		"add":  "vmess.example.com",
		"port": "443",
		"id":   "550e8400-e29b-41d4-a716-446655440000",
		"aid":  "0",
		"scy":  "auto",
		"net":  "ws",
		"type": "none",
		"host": "vmess.example.com",
		"path": "/ws",
		"tls":  "tls",
		"sni":  "vmess.example.com",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	raw := "vmess://" + base64.StdEncoding.EncodeToString(data)

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Protocol != model.ProtocolVMess {
		t.Fatalf("protocol = %q", cfg.Protocol)
	}

	if cfg.Address != "vmess.example.com" {
		t.Fatalf("address = %q", cfg.Address)
	}

	if cfg.Port != 443 {
		t.Fatalf("port = %d", cfg.Port)
	}

	if cfg.UUID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("UUID = %q", cfg.UUID)
	}

	if !cfg.TLS.Enabled {
		t.Fatal("TLS should be enabled")
	}

	if cfg.Transport.Type != "ws" {
		t.Fatalf("transport = %q", cfg.Transport.Type)
	}

	if cfg.Transport.Path != "/ws" {
		t.Fatalf("path = %q", cfg.Transport.Path)
	}
}

func TestParseTrojan(t *testing.T) {
	raw := "trojan://secret-password@example.com:443?sni=example.com&type=ws&path=%2Ftrojan&host=example.com#Trojan-Test"

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Protocol != model.ProtocolTrojan {
		t.Fatalf("protocol = %q", cfg.Protocol)
	}

	if cfg.Address != "example.com" {
		t.Fatalf("address = %q", cfg.Address)
	}

	if cfg.Port != 443 {
		t.Fatalf("port = %d", cfg.Port)
	}

	if cfg.Password != "secret-password" {
		t.Fatalf("password = %q", cfg.Password)
	}

	if !cfg.TLS.Enabled {
		t.Fatal("Trojan TLS should be enabled")
	}

	if cfg.TLS.ServerName != "example.com" {
		t.Fatalf("SNI = %q", cfg.TLS.ServerName)
	}

	if cfg.Transport.Type != "ws" {
		t.Fatalf("transport = %q", cfg.Transport.Type)
	}

	if cfg.Transport.Path != "/trojan" {
		t.Fatalf("path = %q", cfg.Transport.Path)
	}

	if cfg.Name != "Trojan-Test" {
		t.Fatalf("name = %q", cfg.Name)
	}
}

func TestParseUnsupportedProtocol(t *testing.T) {
	_, err := Parse("ss://example")

	if err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestParseEmptyConfig(t *testing.T) {
	_, err := Parse("")

	if err == nil {
		t.Fatal("expected error for empty config")
	}
}
