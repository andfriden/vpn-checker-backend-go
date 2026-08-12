package decoder

import (
	"encoding/base64"
	"strings"
)

func Decode(input string) string {

	data := strings.TrimSpace(input)

	if data == "" {
		return ""
	}

	// Уже готовый URI
	if strings.Contains(data, "://") {
		return data
	}

	// Убираем переносы для Base64
	clean := strings.ReplaceAll(data, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")
	clean = strings.TrimSpace(clean)

	// Обычный Base64
	if decoded, err := base64.StdEncoding.DecodeString(clean); err == nil {
		return string(decoded)
	}

	// URL-safe Base64
	if decoded, err := base64.RawURLEncoding.DecodeString(clean); err == nil {
		return string(decoded)
	}

	return data
}
