package checker

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/andfriden/vpn-checker-backend-go/internal/config"
	"github.com/andfriden/vpn-checker-backend-go/internal/model"
	"github.com/andfriden/vpn-checker-backend-go/internal/parser"
	"github.com/andfriden/vpn-checker-backend-go/internal/singbox"
)

type Checker struct {
	cfg config.CheckerConfig

	singboxBinary   string
	startupTimeout  time.Duration
	shutdownTimeout time.Duration
}

type Result struct {
	Config *model.VPNConfig

	Success bool
	IP      string

	Latency time.Duration

	Error string
}

func New(
	checkerConfig config.CheckerConfig,
	singboxConfig config.SingBoxConfig,
) *Checker {
	startupTimeout := singboxConfig.StartupTimeout
	if startupTimeout <= 0 {
		startupTimeout = 5 * time.Second
	}

	shutdownTimeout := singboxConfig.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 3 * time.Second
	}

	return &Checker{
		cfg:             checkerConfig,
		singboxBinary:   singboxConfig.Binary,
		startupTimeout:  startupTimeout,
		shutdownTimeout: shutdownTimeout,
	}
}

func (c *Checker) Check(
	ctx context.Context,
	raw string,
) *Result {
	result := &Result{}

	raw = strings.TrimSpace(raw)

	if raw == "" {
		result.Error = "empty VPN config"
		return result
	}

	cfg, err := parser.Parse(raw)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Config = cfg

	port, releasePort, err := reserveFreePort()
	if err != nil {
		result.Error = fmt.Sprintf(
			"reserve free SOCKS port: %v",
			err,
		)

		return result
	}

	process, err := singbox.Start(
		ctx,
		c.singboxBinary,
		cfg,
		port,
	)

	releasePort()

	if err != nil {
		result.Error = err.Error()
		return result
	}

	defer func() {
		_ = process.Stop(c.shutdownTimeout)
	}()

	if err := process.WaitReady(
		ctx,
		c.startupTimeout,
	); err != nil {
		stdout, stderr := process.Logs()

		result.Error = formatProcessError(
			err,
			stdout,
			stderr,
		)

		return result
	}

	ipCheckURL := c.cfg.IPCheckURL

	if ipCheckURL == "" {
		ipCheckURL = "https://api.ipify.org"
	}

	ipResult, err := process.CheckIP(
		ctx,
		ipCheckURL,
		c.cfg.Timeout,
	)

	if err != nil {
		stdout, stderr := process.Logs()

		result.Error = formatProcessError(
			err,
			stdout,
			stderr,
		)

		return result
	}

	result.Success = ipResult.Success
	result.IP = strings.TrimSpace(ipResult.IP)
	result.Latency = ipResult.Latency

	if ipResult.Error != nil {
		stdout, stderr := process.Logs()

		result.Error = formatProcessError(
			ipResult.Error,
			stdout,
			stderr,
		)
	}

	return result
}

func (c *Checker) CheckMany(
	ctx context.Context,
	configs []string,
	progress ...func(int, int, int),
) []*Result {

	if len(configs) == 0 {
		return []*Result{}
	}

	var onProgress func(int, int, int)

	if len(progress) > 0 {
		onProgress = progress[0]
	}

	workers := c.cfg.Workers

	if workers < 1 {
		workers = 1
	}

	if workers > len(configs) {
		workers = len(configs)
	}

	maxConcurrent := c.cfg.MaxConcurrentChecks

	if maxConcurrent < 1 {
		maxConcurrent = 1
	}

	if workers > maxConcurrent {
		workers = maxConcurrent
	}

	results := make([]*Result, len(configs))

	jobs := make(chan int)

	var wg sync.WaitGroup

	var progressMu sync.Mutex

	checked := 0
	working := 0
	failed := 0

	wg.Add(workers)

	for workerID := 0; workerID < workers; workerID++ {

		go func(id int) {

			defer wg.Done()

			for {

				select {

				case <-ctx.Done():
					return

				case index, ok := <-jobs:

					if !ok {
						return
					}

					result := c.Check(
						ctx,
						configs[index],
					)

					results[index] = result

					progressMu.Lock()

					checked++

					if result != nil && result.Success {
						working++
					} else {
						failed++
					}

					currentChecked := checked
					currentWorking := working
					currentFailed := failed

					progressMu.Unlock()

					if onProgress != nil {
						onProgress(
							currentChecked,
							currentWorking,
							currentFailed,
						)
					}
				}
			}

		}(workerID)
	}

sendJobs:

	for index := range configs {

		select {

		case <-ctx.Done():
			break sendJobs

		case jobs <- index:
		}
	}

	close(jobs)

	wg.Wait()

	completed := 0

	for _, result := range results {
		if result != nil {
			completed++
		}
	}

	fmt.Println(
		"CheckMany completed:",
		completed,
		"of",
		len(configs),
	)

	for index := range results {

		if results[index] != nil {
			continue
		}

		err := ctx.Err()

		if err == nil {
			err = context.Canceled
		}

		results[index] = &Result{
			Error: err.Error(),
		}
	}

	return results
}

func reserveFreePort() (int, func(), error) {
	listener, err := net.Listen(
		"tcp",
		"127.0.0.1:0",
	)
	if err != nil {
		return 0, func() {}, err
	}

	address, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		_ = listener.Close()

		return 0, func() {}, fmt.Errorf(
			"unexpected listener address type",
		)
	}

	port := address.Port

	release := func() {
		_ = listener.Close()
	}

	return port, release, nil
}

func formatProcessError(
	err error,
	stdout string,
	stderr string,
) string {
	message := err.Error()

	stderr = strings.TrimSpace(stderr)
	stdout = strings.TrimSpace(stdout)

	if stderr != "" {
		message += fmt.Sprintf(
			"; sing-box stderr: %s",
			stderr,
		)
	}

	if stdout != "" {
		message += fmt.Sprintf(
			"; sing-box stdout: %s",
			stdout,
		)
	}

	return message
}
