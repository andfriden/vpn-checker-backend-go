package integration

import (
	"testing"

	"github.com/andfriden/vpn-checker-backend-go/internal/pipeline"
)

func TestRealSourcesPipeline(t *testing.T) {

	urls := []string{
		"https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/refs/heads/main/WHITE-SNI-RU-all.txt",
		"https://raw.githubusercontent.com/kort0881/vpn-vless-configs-russia/refs/heads/main/archive/subscriptions/all.txt",
		"https://raw.githubusercontent.com/zieng2/wl/main/vless.txt",
	}

	result, err := pipeline.Run(urls)

	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	t.Logf(
		"Found configs: %d",
		len(result),
	)

	if len(result) == 0 {
		t.Fatal("no configs found")
	}
}
