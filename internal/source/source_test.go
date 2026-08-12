package source

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {

	content := `
sources:
  - https://test1.com
  - https://test2.com
`

	err := os.WriteFile(
		"test.yaml",
		[]byte(content),
		0644,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove("test.yaml")

	result, err := Load("test.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 sources got %d",
			len(result),
		)
	}
}
