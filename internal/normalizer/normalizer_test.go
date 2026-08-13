package normalizer

import "testing"

func TestUnique(t *testing.T) {
	input := []string{
		"vless://11111111-1111-1111-1111-111111111111@Example.COM:443?security=tls#Germany",
		"vless://11111111-1111-1111-1111-111111111111@example.com:443?security=tls#Frankfurt",

		"vless://22222222-2222-2222-2222-222222222222@example.com:443?security=tls",

		"hy2://password@example.com:443#one",
		"hysteria2://password@example.com:443#two",
	}

	result := Unique(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 unique configs, got %d", len(result))
	}
}
