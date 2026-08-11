package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/checker"
	"github.com/andfriden/vpn-checker-backend-go/internal/config"
	"github.com/andfriden/vpn-checker-backend-go/internal/downloader"
	"github.com/andfriden/vpn-checker-backend-go/internal/model"
	"github.com/andfriden/vpn-checker-backend-go/internal/storage"
)

type Status struct {
	Status   string `json:"status"`
	Total    int    `json:"total"`
	Checked  int    `json:"checked"`
	Working  int    `json:"working"`
	Failed   int    `json:"failed"`
	Progress int    `json:"progress"`
}

type Runner struct {
	cfg     *config.Config
	checker *checker.Checker
	storage *storage.FileStorage

	mu      sync.Mutex
	running bool
	status  Status
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

		status: Status{
			Status: "idle",
		},
	}
}

func (r *Runner) IsRunning() bool {

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.running
}

func (r *Runner) Status() Status {

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.status
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

	r.status = Status{
		Status: "running",
	}

	r.mu.Unlock()

	go func() {

		defer func() {

			r.mu.Lock()
			r.running = false

			if r.status.Progress == 100 {
				r.status.Status = "completed"
			}

			r.mu.Unlock()

		}()

		if err := r.Run(ctx, url); err != nil {

			r.mu.Lock()

			r.status.Status = "error"

			r.mu.Unlock()

			fmt.Println(
				"check failed:",
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

	configs, err := downloader.DownloadURL(url)

	if err != nil {
		return err
	}

	total := len(configs)

	r.mu.Lock()

	r.status.Total = total

	r.mu.Unlock()

	results := r.checker.CheckMany(
		ctx,
		configs,
		func(
			checked int,
			working int,
			failed int,
		) {

			progress := 0

			if total > 0 {
				progress = checked * 100 / total
			}

			r.mu.Lock()

			r.status.Status = "running"
			r.status.Checked = checked
			r.status.Working = working
			r.status.Failed = failed
			r.status.Progress = progress

			r.mu.Unlock()

		},
	)

	fmt.Printf(
		"Checked: %d\n",
		len(results),
	)

	storageResults :=
		make(
			[]model.CheckResult,
			0,
			len(results),
		)

	workingConfigs :=
		make(
			[]*model.VPNConfig,
			0,
		)

	protocolStats :=
		make(
			map[string]model.ProtocolStats,
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

		storageResults =
			append(
				storageResults,
				checkResult,
			)

		if result.Config != nil {

			protocol :=
				string(result.Config.Protocol)

			stat :=
				protocolStats[protocol]

			stat.Total++

			if result.Success {

				stat.Working++

			} else {

				stat.Failed++
			}

			protocolStats[protocol] = stat
		}

		if result.Success {

			working++

			if result.Config != nil {

				workingConfigs =
					append(
						workingConfigs,
						result.Config,
					)
			}
		}
	}

	if err := r.storage.SaveResults(
		storageResults,
	); err != nil {

		return err
	}

	if err := r.storage.ExportWorking(
		workingConfigs,
	); err != nil {

		return err
	}

	stats := model.Stats{

		Total: len(storageResults),

		Checked: len(storageResults),

		Working: working,

		Failed: len(storageResults) - working,

		Protocols: protocolStats,

		UpdatedAt: time.Now(),
	}

	if err := r.storage.SaveStats(
		stats,
	); err != nil {

		return err
	}

	r.mu.Lock()

	r.status.Status = "completed"

	r.status.Progress = 100

	r.mu.Unlock()

	fmt.Printf(
		"Working VPN: %d\n",
		working,
	)

	return nil
}
