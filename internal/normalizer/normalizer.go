package normalizer

import (
	"net/url"
	"strings"
)

func Unique(configs []string) []string {

	seen := make(map[string]bool)
	result := make([]string, 0, len(configs))

	for _, cfg := range configs {

		cfg = strings.TrimSpace(cfg)

		if cfg == "" {
			continue
		}

		key := normalizeKey(cfg)

		if seen[key] {
			continue
		}

		seen[key] = true
		result = append(result, cfg)
	}

	return result
}

func normalizeKey(cfg string) string {

	u, err := url.Parse(cfg)

	if err != nil {
		return cfg
	}

	host := u.Hostname()
	port := u.Port()

	scheme := u.Scheme

	switch scheme {

	case "vless",
		"vmess",
		"trojan",
		"hysteria",
		"hysteria2",
		"hy2":

		return strings.Join(
			[]string{
				scheme,
				host,
				port,
				u.User.String(),
			},
			"|",
		)

	case "ss":

		return strings.Join(
			[]string{
				scheme,
				host,
				port,
				u.User.String(),
			},
			"|",
		)

	default:

		return cfg
	}
}
