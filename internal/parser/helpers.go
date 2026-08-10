package parser

import "net/url"

func fragmentName(fragment string) string {

	if fragment == "" {
		return ""
	}

	name, err := url.QueryUnescape(fragment)

	if err != nil {
		return fragment
	}

	return name
}
