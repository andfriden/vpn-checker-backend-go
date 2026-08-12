package normalizer

import "testing"

func TestUnique(t *testing.T) {

	input := []string{

		"vless://111@server.com:443?type=tcp#Europe",
		"vless://111@server.com:443?type=tcp#Germany",
		"vless://222@server.com:443?type=tcp#USA",
	}

	result := Unique(input)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 configs, got %d",
			len(result),
		)
	}
}
