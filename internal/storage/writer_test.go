package storage

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {

	file := "testdata/configs.txt"

	configs := []string{
		"vless://test1",
		"ss://test2",
	}

	err := Save(file, configs)

	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll("testdata")

	data, err := os.ReadFile(file)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("empty file")
	}
}
