package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/app"
)

type Scheduler struct {
	runner   *app.Runner
	interval time.Duration
	url      string
}

func New(
	runner *app.Runner,
	interval time.Duration,
	url string,
) *Scheduler {

	if interval <= 0 {
		interval = 30 * time.Minute
	}

	return &Scheduler{
		runner:   runner,
		interval: interval,
		url:      url,
	}
}

func (s *Scheduler) Start(
	ctx context.Context,
) {

	log.Printf(
		"scheduler started, interval: %s",
		s.interval,
	)

	// первый запуск сразу
	if !s.runner.RunAsync(
		ctx,
		s.url,
	) {
		log.Println(
			"initial check skipped: already running",
		)
	}

	ticker := time.NewTicker(
		s.interval,
	)

	defer ticker.Stop()

	for {

		select {

		case <-ctx.Done():

			log.Println(
				"scheduler stopped",
			)

			return

		case <-ticker.C:

			log.Println(
				"scheduled VPN check started",
			)

			if !s.runner.RunAsync(
				ctx,
				s.url,
			) {

				log.Println(
					"scheduled VPN check skipped: already running",
				)

				continue
			}

			log.Println(
				"scheduled VPN check launched",
			)
		}
	}
}
