package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/api"
	"github.com/andfriden/vpn-checker-backend-go/internal/app"
	"github.com/andfriden/vpn-checker-backend-go/internal/config"
	"github.com/andfriden/vpn-checker-backend-go/internal/scheduler"
	"github.com/andfriden/vpn-checker-backend-go/internal/service"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

func main() {

	configPath := "configs/config.yaml"

	if v := os.Getenv("VPN_CHECKER_CONFIG"); v != "" {
		configPath = v
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	runner := app.New(cfg)

	fileStorage := storage.New(
		cfg.Storage.Path,
	)

	resultsService := service.NewResultsService(
		fileStorage,
	)

	blackURL := "https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/main/BLACK_SS%2BAll_RUS.txt"

	whiteURL := "https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/main/WHITE_SS%2BAll_RUS.txt"

	go scheduler.New(
		runner,
		cfg.Checker.HealthCheckInterval,
		blackURL,
		whiteURL,
	).Start(ctx)
	handler := api.NewHandler(
		runner,
		resultsService,
	)
	server := &http.Server{

		Addr: fmt.Sprintf(
			"%s:%d",
			cfg.Server.Host,
			cfg.Server.Port,
		),

		Handler: api.Routes(
			handler,
		),

		ReadTimeout: cfg.Server.ReadTimeout,

		WriteTimeout: cfg.Server.WriteTimeout,

		IdleTimeout: cfg.Server.IdleTimeout,
	}

	go func() {

		log.Printf(
			"VPN Checker API listening on %s",
			server.Addr,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server error: %v",
				err,
			)
		}

	}()

	<-ctx.Done()

	log.Println(
		"shutdown signal received",
	)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(
		shutdownCtx,
	); err != nil {

		log.Printf(
			"server shutdown error: %v",
			err,
		)
	}

	log.Println(
		"VPN Checker stopped",
	)
}
