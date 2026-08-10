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

type Runner struct {
	cfg     *config.Config
	checker *checker.Checker
	storage *storage.FileStorage

	mu      sync.RWMutex
	running bool
	status  string
	total   int
	checked int
	working int
	failed  int
	err     string
}

type Status struct {
	Status   string `json:"status"`
	Total    int    `json:"total"`
	Checked  int    `json:"checked"`
	Working  int    `json:"working"`
	Failed   int    `json:"failed"`
	Progress int    `json:"progress"`
	Error    string `json:"error,omitempty"`
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
		status: "idle",
	}
}

func (r *Runner) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.running
}

func (r *Runner) Status() Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	progress := 0

	if r.total > 0 {
		progress = r.checked * 100 / r.total
	}

	return Status{
		Status:   r.status,
		Total:    r.total,
		Checked:  r.checked,
		Working:  r.working,
		Failed:   r.failed,
		Progress: progress,
		Error:    r.err,
	}
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
	r.status = "running"
	r.total = 0
	r.checked = 0
	r.working = 0
	r.failed = 0
	r.err = ""

	r.mu.Unlock()

	go func() {
		defer func() {
			r.mu.Lock()
			r.running = false

			if r.err != "" {
				r.status = "error"
			} else {
				r.status = "completed"
			}

			r.mu.Unlock()
		}()

		if err := r.Run(ctx, url); err != nil {
			fmt.Printf(
				"check failed: %v\n",
				err,
			)

			r.mu.Lock()
			r.err = err.Error()
			r.mu.Unlock()
		}
	}()

	return true
}

func (r *Runner) Run(
	ctx context.Context,
	url string,
) error {
	fmt.Println("Downloading configs...")

	configs, err := downloader.DownloadURL(url)
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

	r.mu.Lock()
	r.total = len(configs)
	r.mu.Unlock()

	fmt.Println("Checking VPN configs...")

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

	workingConfigs := make(
		[]*model.VPNConfig,
		0,
	)

	protocolStats := make(
		map[string]model.ProtocolStats,
	)

	working := 0
	checked := 0

	for _, result := range results {
		if result == nil {
			continue
		}

		checked++

		checkResult := model.CheckResult{
			Config:  result.Config,
			Success: result.Success,
			IP:      result.IP,
			Latency: result.Latency,
			Error:   result.Error,
		}

		storageResults = append(
			storageResults,
			checkResult,
		)

		if result.Config != nil {
			protocol := string(
				result.Config.Protocol,
			)

			protocolResult := protocolStats[protocol]

			protocolResult.Total++

			if result.Success {
				protocolResult.Working++
			} else {
				protocolResult.Failed++
			}

			protocolStats[protocol] = protocolResult
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

		r.mu.Lock()
		r.checked = checked
		r.working = working
		r.failed = checked - working
		r.mu.Unlock()
	}

	failed := len(storageResults) - working

	if err := r.storage.SaveResults(
		storageResults,
	); err != nil {
		return fmt.Errorf(
			"save results: %w",
			err,
		)
	}

	if err := r.storage.ExportWorking(
		workingConfigs,
	); err != nil {
		return fmt.Errorf(
			"export working configs: %w",
			err,
		)
	}

	stats := model.Stats{
		Total:     len(storageResults),
		Checked:   len(storageResults),
		Working:   working,
		Failed:    failed,
		Protocols: protocolStats,
		UpdatedAt: time.Now(),
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
		"Working configs exported: %s/all-working.txt\n",
		r.cfg.Storage.Path,
	)

	fmt.Printf(
		"Stats saved: %s/stats.json\n",
		r.cfg.Storage.Path,
	)

	return nil
}
