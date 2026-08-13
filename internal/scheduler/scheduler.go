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
}

func New(
	runner *app.Runner,
	interval time.Duration,
) *Scheduler {

	if interval <= 0 {
		interval = 30 * time.Minute
	}

	return &Scheduler{
		runner:   runner,
		interval: interval,
	}
}

func (s *Scheduler) Start(
	ctx context.Context,
) {

	log.Printf(
		"scheduler started, interval: %s",
		s.interval,
	)

	// Первый запуск сразу.
	s.run(ctx)

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

			s.run(ctx)
		}
	}
}

func (s *Scheduler) run(
	ctx context.Context,
) {

	if !s.runner.RunCollectedAsync(ctx) {
		log.Println(
			"scheduled check skipped: already running",
		)
	}
}
