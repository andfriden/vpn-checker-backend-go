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

	blackURL string
	whiteURL string
}

func New(
	runner *app.Runner,
	interval time.Duration,
	blackURL string,
	whiteURL string,
) *Scheduler {

	if interval <= 0 {
		interval = 30 * time.Minute
	}

	return &Scheduler{
		runner:   runner,
		interval: interval,

		blackURL: blackURL,
		whiteURL: whiteURL,
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

	// BLACK
	if s.blackURL != "" {

		log.Println(
			"starting BLACK list check",
		)

		if !s.runner.RunAsync(
			ctx,
			s.blackURL,
		) {

			log.Println(
				"BLACK check skipped: already running",
			)
		}
	}

	// WHITE
	if s.whiteURL != "" {

		log.Println(
			"starting WHITE list check",
		)

		if !s.runner.RunAsync(
			ctx,
			s.whiteURL,
		) {

			log.Println(
				"WHITE check skipped: already running",
			)
		}
	}
}
