package main

import (
	"log"

	"github.com/andfriden/vpn-checker-backend-go/internal/pipeline"
	"github.com/andfriden/vpn-checker-backend-go/internal/source"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

func main() {

	urls, err := source.Load(
		"configs/sources.yaml",
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"Loaded sources: %d",
		len(urls),
	)

	configs, err := pipeline.Run(urls)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"Collected configs: %d",
		len(configs),
	)

	err = storage.Save(
		"data/configs/all.txt",
		configs,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		"Saved to data/configs/all.txt",
	)
}
