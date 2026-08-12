package decoder

import "testing"

func TestDecodeBase64(t *testing.T) {

	input := "dmxlc3M6Ly90ZXN0"

	result := Decode(input)

	if result != "vless://test" {
		t.Fatalf(
			"expected vless://test, got %s",
			result,
		)
	}
}

func TestDecodePlainConfig(t *testing.T) {

	input := "vless://example"

	result := Decode(input)

	if result != "vless://example" {
		t.Fatalf(
			"expected vless://example, got %s",
			result,
		)
	}
}

func TestDecodeBase64WithNewLines(t *testing.T) {

	input := `
dmxlc3M6Ly90ZXN0
`

	result := Decode(input)

	if result != "vless://test" {
		t.Fatalf(
			"expected vless://test, got %s",
			result,
		)
	}
}

func TestDecodeInvalidBase64(t *testing.T) {

	input := "random text"

	result := Decode(input)

	if result != "random text" {
		t.Fatalf(
			"expected original text, got %s",
			result,
		)
	}
}
