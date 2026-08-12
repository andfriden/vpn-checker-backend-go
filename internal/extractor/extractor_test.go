package extractor

import "testing"

func TestExtract(t *testing.T) {
	input := `
some random text

vless://test-config-1

another line

trojan://test-config-2
ss://test-config-3

not a config
`

	result := Extract(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 configs, got %d", len(result))
	}
}
