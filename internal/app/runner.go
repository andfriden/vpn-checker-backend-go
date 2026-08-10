package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/andfriden/vpn-checker-backend-go/internal/checker"
	"github.com/andfriden/vpn-checker-backend-go/internal/config"
	"github.com/andfriden/vpn-checker-backend-go/internal/downloader"
	"github.com/andfriden/vpn-checker-backend-go/internal/model"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

type Runner struct {
	cfg     *config.Config
	checker *checker.Checker
	storage *storage.FileStorage

	mu      sync.Mutex
	running bool
}

func New(cfg *config.Config) *Runner {

	return &Runner{

		cfg: cfg,

		checker: checker.New(
			cfg.Checker,
			cfg.SingBox,
		),

		storage: storage.New(
			cfg.Storage.Path,
		),
	}
}

func (r *Runner) IsRunning() bool {

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.running
}

func (r *Runner) RunAsync(
	ctx context.Context,
	url string,
) bool {

	r.mu.Lock()

	if r.running {

		r.mu.Unlock()

		return false
	}

	r.running = true

	r.mu.Unlock()

	go func() {

		defer func() {

			r.mu.Lock()

			r.running = false

			r.mu.Unlock()

		}()

		if err := r.Run(
			ctx,
			url,
		); err != nil {

			fmt.Printf(
				"check failed: %v\n",
				err,
			)

		}

	}()

	return true
}

func (r *Runner) Run(
	ctx context.Context,
	url string,
) error {

	fmt.Println(
		"Downloading configs...",
	)

	configs, err := downloader.DownloadURL(
		url,
	)

	if err != nil {

		return fmt.Errorf(
			"download configs: %w",
			err,
		)
	}

	if len(configs) == 0 {

		return fmt.Errorf(
			"no VPN configs downloaded",
		)
	}

	fmt.Printf(
		"Loaded configs: %d\n",
		len(configs),
	)

	fmt.Println(
		"Checking VPN configs...",
	)

	results := r.checker.CheckMany(
		ctx,
		configs,
	)

	fmt.Printf(
		"Checked: %d\n",
		len(results),
	)

	storageResults := make(
		[]model.CheckResult,
		0,
		len(results),
	)

	working := 0

	for _, result := range results {

		if result == nil {

			continue
		}

		checkResult := model.CheckResult{

			Config: result.Config,

			Success: result.Success,

			IP: result.IP,

			Latency: result.Latency,

			Error: result.Error,
		}

		storageResults = append(
			storageResults,
			checkResult,
		)

		if result.Success {

			working++

		}

	}

	if err := r.storage.SaveResults(
		storageResults,
	); err != nil {

		return fmt.Errorf(
			"save results: %w",
			err,
		)
	}

	stats := model.Stats{

		Total: len(storageResults),

		Checked: len(storageResults),

		Working: working,

		Failed: len(storageResults) - working,
	}

	if err := r.storage.SaveStats(
		stats,
	); err != nil {

		return fmt.Errorf(
			"save stats: %w",
			err,
		)
	}

	fmt.Printf(
		"Working VPN: %d\n",
		working,
	)

	fmt.Printf(
		"Results saved: %s/results.json\n",
		r.cfg.Storage.Path,
	)

	fmt.Printf(
		"Stats saved: %s/stats.json\n",
		r.cfg.Storage.Path,
	)

	return nil
}
