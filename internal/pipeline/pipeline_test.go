package pipeline

import "testing"

func TestRun(t *testing.T) {
	urls := []string{}

	result, err := Run(urls)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
