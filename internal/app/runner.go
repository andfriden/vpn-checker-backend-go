package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

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
	source string,
) bool {
	r.mu.Lock()

	if r.running {
		r.mu.Unlock()
		return false
	}

	started := time.Now()

	r.running = true

	r.status = Status{
		Status:    "running",
		StartedAt: &started,
	}

	r.mu.Unlock()

	go func() {
		err := r.Run(
			ctx,
			url,
			source,
		)

		finished := time.Now()

		r.mu.Lock()
		defer r.mu.Unlock()

		r.running = false
		r.status.FinishedAt = &finished

		if r.status.StartedAt != nil {
			r.status.ElapsedSeconds = finished.
				Sub(*r.status.StartedAt).
				Seconds()
		}

		if err != nil {
			r.status.Status = "error"

			fmt.Println(
				"check failed:",
				err,
			)

			return
		}

		r.status.Status = "completed"
		r.status.Checked = r.status.Total
		r.status.Progress = 100
		r.status.EstimatedSecondsLeft = 0
	}()

	return true
}

func (r *Runner) Run(
	ctx context.Context,
	url string,
	source string,
) error {
	fmt.Println(
		"Downloading configs:",
		source,
	)

	configs, err := downloader.DownloadURL(url)
	if err != nil {
		return err
	}

	return r.RunConfigs(
		ctx,
		configs,
		source,
	)
}

func (r *Runner) RunCollectedAsync(
	ctx context.Context,
) bool {
	r.mu.Lock()

	if r.running {
		r.mu.Unlock()
		return false
	}

	started := time.Now()

	r.running = true

	r.status = Status{
		Status:    "running",
		StartedAt: &started,
	}

	r.mu.Unlock()

	go func() {
		err := r.RunCollected(ctx)

		finished := time.Now()

		r.mu.Lock()
		defer r.mu.Unlock()

		r.running = false
		r.status.FinishedAt = &finished

		if r.status.StartedAt != nil {
			r.status.ElapsedSeconds = finished.
				Sub(*r.status.StartedAt).
				Seconds()
		}

		if err != nil {
			r.status.Status = "error"

			fmt.Println(
				"check failed:",
				err,
			)

			return
		}

		r.status.Status = "completed"
		r.status.Progress = 100
		r.status.EstimatedSecondsLeft = 0

		if r.status.Checked > 0 &&
			r.status.ElapsedSeconds > 0 {

			r.status.CurrentSpeed =
				float64(r.status.Checked) /
					r.status.ElapsedSeconds
		}
	}()

	return true
}

func (r *Runner) RunCollected(
	ctx context.Context,
) error {
	data, err := os.ReadFile(
		"data/configs/all.txt",
	)

	if err != nil {
		return err
	}

	lines := strings.Split(
		string(data),
		"\n",
	)

	configs := make(
		[]string,
		0,
		len(lines),
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		configs = append(
			configs,
			line,
		)
	}

	return r.RunConfigs(
		ctx,
		configs,
		"collector",
	)
}

func (r *Runner) RunConfigs(
	ctx context.Context,
	configs []string,
	source string,
) error {
	total := len(configs)

	r.mu.Lock()

	r.status.Total = total
	r.status.Checked = 0
	r.status.Working = 0
	r.status.Failed = 0
	r.status.Progress = 0
	r.status.CurrentSpeed = 0
	r.status.EstimatedSecondsLeft = 0

	r.mu.Unlock()

	results := r.checker.CheckMany(
		ctx,
		configs,
		source,
		func(
			checked int,
			working int,
			failed int,
		) {
			r.mu.Lock()
			defer r.mu.Unlock()

			progress := 0

			if total > 0 {
				progress = checked * 100 / total
			}

			r.status.Checked = checked
			r.status.Working = working
			r.status.Failed = failed
			r.status.Progress = progress

			if r.status.StartedAt != nil {
				elapsed := time.Since(
					*r.status.StartedAt,
				)

				r.status.ElapsedSeconds =
					elapsed.Seconds()

				if elapsed > 0 {
					r.status.CurrentSpeed =
						float64(checked) /
							elapsed.Seconds()

					if r.status.CurrentSpeed > 0 {
						remaining := total - checked

						if remaining < 0 {
							remaining = 0
						}

						r.status.EstimatedSecondsLeft =
							float64(remaining) /
								r.status.CurrentSpeed
					}
				}
			}
		},
	)

	storageResults := make(
		[]model.CheckResult,
		0,
		len(results),
	)

	workingConfigs := make(
		[]*model.VPNConfig,
		0,
	)

	protocolStats := make(
		map[string]model.ProtocolStats,
	)

	working := 0

	for _, result := range results {
		if result == nil {
			continue
		}

		storageResults = append(
			storageResults,
			model.CheckResult{
				Config:  result.Config,
				Success: result.Success,
				IP:      result.IP,
				Latency: result.Latency,
				Error:   result.Error,
			},
		)

		if result.Config != nil {
			protocol := string(
				result.Config.Protocol,
			)

			stat := protocolStats[protocol]

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
				workingConfigs = append(
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
		Total:     len(storageResults),
		Checked:   len(storageResults),
		Working:   working,
		Failed:    len(storageResults) - working,
		Protocols: protocolStats,
		UpdatedAt: time.Now(),
	}

	if err := r.storage.SaveStats(
		stats,
	); err != nil {
		return err
	}

	return nil
}
