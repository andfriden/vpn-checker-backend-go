package downloader

import (
	"io"
	"net/http"
	"time"
)

type Downloader struct {
	Client *http.Client
}

func New() *Downloader {
	return &Downloader{
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (d *Downloader) Fetch(url string) (string, error) {

	resp, err := d.Client.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
