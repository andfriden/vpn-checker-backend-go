package downloader

import "testing"

func TestDownloader(t *testing.T) {

	d := New()

	data, err := d.Fetch(
		"https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/refs/heads/main/WHITE-SNI-RU-all.txt",
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("empty response")
	}

	t.Logf("Downloaded: %d bytes", len(data))
}
