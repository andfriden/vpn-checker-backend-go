package normalizer

import (
	"net/url"
	"strconv"
	"strings"
)

func Unique(configs []string) []string {
	seen := make(map[string]struct{}, len(configs))

	result := make([]string, 0, len(configs))

	for _, cfg := range configs {
		cfg = strings.TrimSpace(cfg)

		if cfg == "" {
			continue
		}

		key := normalizeKey(cfg)

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, cfg)
	}

	return result
}

func normalizeKey(cfg string) string {
	u, err := url.Parse(cfg)
	if err != nil {
		return cfg
	}

	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	host := strings.ToLower(strings.TrimSpace(u.Hostname()))
	port := normalizePort(scheme, u.Port())
	user := ""

	if u.User != nil {
		user = u.User.String()
	}

	switch scheme {
	case "vless",
		"vmess",
		"trojan",
		"hysteria",
		"hysteria2",
		"hy2":

		if scheme == "hy2" {
			scheme = "hysteria2"
		}

		return strings.Join(
			[]string{
				scheme,
				host,
				port,
				user,
			},
			"|",
		)

	case "ss":
		return strings.Join(
			[]string{
				scheme,
				host,
				port,
				user,
			},
			"|",
		)

	default:
		return cfg
	}
}

func normalizePort(scheme string, port string) string {
	if port != "" {
		return port
	}

	switch scheme {
	case "vless",
		"vmess",
		"trojan",
		"hysteria",
		"hysteria2",
		"hy2",
		"ss":
		return strconv.Itoa(defaultPort(scheme))
	default:
		return ""
	}
}

func defaultPort(scheme string) int {
	switch scheme {
	case "hysteria", "hysteria2", "hy2":
		return 443
	default:
		return 443
	}
}
