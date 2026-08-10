package checker

import (
	"context"
	"testing"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/config"
)

func testChecker() *Checker {
	return New(
		config.CheckerConfig{
			Workers:             2,
			MaxConcurrentChecks: 2,
			Timeout:             5 * time.Second,
			IPCheckURL:          "https://api.ipify.org",
		},
		config.SingBoxConfig{
			Binary:          "sing-box",
			StartupTimeout:  2 * time.Second,
			ShutdownTimeout: 1 * time.Second,
		},
	)
}

func TestCheckerInvalidConfig(t *testing.T) {
	checker := testChecker()

	result := checker.Check(
		context.Background(),
		"invalid://vpn-config",
	)

	if result == nil {
		t.Fatal("result is nil")
	}

	if result.Success {
		t.Fatal("invalid config cannot be successful")
	}

	if result.Error == "" {
		t.Fatal("expected parser error")
	}
}

func TestReserveFreePort(t *testing.T) {
	port, release, err := reserveFreePort()

	if err != nil {
		t.Fatalf(
			"reserveFreePort() error = %v",
			err,
		)
	}

	if port < 1 || port > 65535 {
		t.Fatalf(
			"invalid port: %d",
			port,
		)
	}

	if release == nil {
		t.Fatal("release function is nil")
	}

	release()
}

func TestCheckManyEmpty(t *testing.T) {
	checker := testChecker()

	results := checker.CheckMany(
		context.Background(),
		nil,
	)

	if results == nil {
		t.Fatal("results should not be nil")
	}

	if len(results) != 0 {
		t.Fatalf(
			"results length = %d, want 0",
			len(results),
		)
	}
}

func TestCheckManyInvalidConfigs(t *testing.T) {
	checker := testChecker()

	configs := []string{
		"invalid://one",
		"invalid://two",
		"invalid://three",
	}

	results := checker.CheckMany(
		context.Background(),
		configs,
	)

	if results == nil {
		t.Fatal("results should not be nil")
	}

	if len(results) != len(configs) {
		t.Fatalf(
			"results length = %d, want %d",
			len(results),
			len(configs),
		)
	}

	for index, result := range results {
		if result == nil {
			t.Fatalf(
				"result %d is nil",
				index,
			)
		}

		if result.Success {
			t.Fatalf(
				"result %d unexpectedly successful",
				index,
			)
		}

		if result.Error == "" {
			t.Fatalf(
				"result %d has no error",
				index,
			)
		}
	}
}

func TestCheckManyContextCancellation(t *testing.T) {
	checker := testChecker()

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	cancel()

	configs := []string{
		"invalid://one",
		"invalid://two",
		"invalid://three",
		"invalid://four",
	}

	results := checker.CheckMany(
		ctx,
		configs,
	)

	if results == nil {
		t.Fatal("results should not be nil")
	}

	if len(results) != len(configs) {
		t.Fatalf(
			"results length = %d, want %d",
			len(results),
			len(configs),
		)
	}

	for index, result := range results {
		if result == nil {
			t.Fatalf(
				"result %d is nil",
				index,
			)
		}

		if result.Error == "" {
			t.Fatalf(
				"result %d should contain cancellation error",
				index,
			)
		}
	}
}
