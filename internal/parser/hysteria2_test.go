package parser

import (
	"testing"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

func TestParseHysteria2(t *testing.T) {
	raw := "hysteria2://secret-password@example.com:443?sni=example.com&insecure=0#Hysteria2-Test"

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Protocol != model.ProtocolHysteria2 {
		t.Fatalf("protocol = %q, want %q", cfg.Protocol, model.ProtocolHysteria2)
	}

	if cfg.Address != "example.com" {
		t.Fatalf("address = %q, want %q", cfg.Address, "example.com")
	}

	if cfg.Port != 443 {
		t.Fatalf("port = %d, want %d", cfg.Port, 443)
	}

	if cfg.Password != "secret-password" {
		t.Fatalf("password = %q", cfg.Password)
	}

	if !cfg.TLS.Enabled {
		t.Fatal("TLS should be enabled")
	}

	if cfg.TLS.ServerName != "example.com" {
		t.Fatalf("SNI = %q", cfg.TLS.ServerName)
	}

	if cfg.Protocol != model.ProtocolHysteria2 {
		t.Fatalf("protocol = %q, want hysteria2",
			cfg.Protocol)
	}

	if cfg.Name != "Hysteria2-Test" {
		t.Fatalf("name = %q", cfg.Name)
	}
}

func TestParseHy2(t *testing.T) {
	raw := "hy2://password123@server.example:8443?sni=server.example#HY2-Test"

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.Protocol != model.ProtocolHysteria2 {
		t.Fatalf("protocol = %q", cfg.Protocol)
	}

	if cfg.Address != "server.example" {
		t.Fatalf("address = %q", cfg.Address)
	}

	if cfg.Port != 8443 {
		t.Fatalf("port = %d", cfg.Port)
	}

	if cfg.Password != "password123" {
		t.Fatalf("password = %q", cfg.Password)
	}
}

func TestParseHysteria2Insecure(t *testing.T) {
	raw := "hysteria2://password@example.com:443?sni=example.com&insecure=1"

	cfg, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if !cfg.TLS.Enabled {
		t.Fatal("TLS should remain enabled for insecure=1")
	}

	if !cfg.TLS.Insecure {
		t.Fatal("TLS verification should be disabled for insecure=1")
	}
}
