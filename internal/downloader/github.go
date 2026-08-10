package downloader

import (
	"fmt"
	"io"
	"net/http"
)

func DownloadURL(
	url string,
) ([]string, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf(
			"download failed: %s",
			resp.Status,
		)
	}

	body, err := io.ReadAll(
		resp.Body,
	)

	if err != nil {
		return nil, err
	}

	return New().ParseList(
		string(body),
	), nil
}
