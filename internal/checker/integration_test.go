package checker

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/config"
)

func TestIntegrationSingleVPN(t *testing.T) {
	if os.Getenv("VPN_CHECKER_INTEGRATION") != "1" {
		t.Skip(
			"set VPN_CHECKER_INTEGRATION=1 to run integration test",
		)
	}

	raw := os.Getenv("VPN_CHECKER_TEST_CONFIG")

	if raw == "" {
		t.Fatal(
			"VPN_CHECKER_TEST_CONFIG is not set",
		)
	}

	binary := os.Getenv("VPN_CHECKER_SINGBOX")

	if binary == "" {
		binary = "sing-box"
	}

	checker := New(
		config.CheckerConfig{
			Workers:             1,
			MaxConcurrentChecks: 1,
			Timeout:             15 * time.Second,
			IPCheckURL:          "https://api.ipify.org",
		},
		config.SingBoxConfig{
			Binary:          binary,
			StartupTimeout:  10 * time.Second,
			ShutdownTimeout: 3 * time.Second,
		},
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	result := checker.Check(
		ctx,
		raw,
	)

	if result == nil {
		t.Fatal("result is nil")
	}

	t.Logf(
		"success: %v",
		result.Success,
	)

	t.Logf(
		"IP: %s",
		result.IP,
	)

	t.Logf(
		"latency: %s",
		result.Latency,
	)

	t.Logf(
		"error: %s",
		result.Error,
	)

	if result.Config != nil {
		t.Logf(
			"protocol: %s",
			result.Config.Protocol,
		)

		t.Logf(
			"server: %s:%d",
			result.Config.Address,
			result.Config.Port,
		)
		t.Logf(
			"FULL CONFIG: %+v",
			result.Config,
		)

	}

	if !result.Success {
		t.Fatalf(
			"VPN check failed: %s",
			result.Error,
		)
	}
}
