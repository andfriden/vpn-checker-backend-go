package singbox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/net/proxy"

	"github.com/andfriden/vpn-checker-backend-go/internal/model"
)

type Process struct {
	cmd        *exec.Cmd
	tempDir    string
	configPath string
	socksPort  int
	done       chan error
}

type CheckResult struct {
	Success bool
	IP      string
	Latency time.Duration
	Error   error
}

func Start(
	ctx context.Context,
	binary string,
	vpnConfig *model.VPNConfig,
	socksPort int,
) (*Process, error) {
	if vpnConfig == nil {
		return nil, fmt.Errorf("VPN config is nil")
	}

	if binary == "" {
		return nil, fmt.Errorf("sing-box binary is empty")
	}

	if socksPort < 1 || socksPort > 65535 {
		return nil, fmt.Errorf(
			"invalid SOCKS port: %d",
			socksPort,
		)
	}

	config, err := Build(vpnConfig, socksPort)
	if err != nil {
		return nil, fmt.Errorf(
			"build sing-box config: %w",
			err,
		)
	}

	data, err := json.MarshalIndent(
		config,
		"",
		"  ",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal sing-box config: %w",
			err,
		)
	}

	tempDir, err := os.MkdirTemp(
		"",
		"vpn-checker-*",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create temporary directory: %w",
			err,
		)
	}

	configPath := filepath.Join(
		tempDir,
		"config.json",
	)

	if err := os.WriteFile(
		configPath,
		data,
		0600,
	); err != nil {
		_ = os.RemoveAll(tempDir)

		return nil, fmt.Errorf(
			"write sing-box config: %w",
			err,
		)
	}

	// Temporary diagnostic output.
	// We will remove this after the Hysteria2
	// configuration is confirmed to work.
	fmt.Printf(
		"\n===== SING-BOX CONFIG =====\n%s\n===== END CONFIG =====\n\n",
		string(data),
	)

	stdoutPath := filepath.Join(
		tempDir,
		"stdout.log",
	)

	stderrPath := filepath.Join(
		tempDir,
		"stderr.log",
	)

	stdout, err := os.OpenFile(
		stdoutPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	if err != nil {
		_ = os.RemoveAll(tempDir)

		return nil, fmt.Errorf(
			"create stdout log: %w",
			err,
		)
	}

	stderr, err := os.OpenFile(
		stderrPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	if err != nil {
		_ = stdout.Close()
		_ = os.RemoveAll(tempDir)

		return nil, fmt.Errorf(
			"create stderr log: %w",
			err,
		)
	}

	cmd := exec.CommandContext(
		ctx,
		binary,
		"run",
		"-c",
		configPath,
	)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		_ = os.RemoveAll(tempDir)

		return nil, fmt.Errorf(
			"start sing-box: %w",
			err,
		)
	}

	done := make(chan error, 1)

	process := &Process{
		cmd:        cmd,
		tempDir:    tempDir,
		configPath: configPath,
		socksPort:  socksPort,
		done:       done,
	}

	go func() {
		err := cmd.Wait()

		_ = stdout.Close()
		_ = stderr.Close()

		done <- err
		close(done)
	}()

	return process, nil
}

func (p *Process) WaitReady(
	ctx context.Context,
	timeout time.Duration,
) error {
	if p == nil {
		return fmt.Errorf("process is nil")
	}

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	address := fmt.Sprintf(
		"127.0.0.1:%d",
		p.socksPort,
	)

	waitCtx, cancel := context.WithTimeout(
		ctx,
		timeout,
	)
	defer cancel()

	ticker := time.NewTicker(
		50 * time.Millisecond,
	)
	defer ticker.Stop()

	for {
		select {
		case <-waitCtx.Done():
			return fmt.Errorf(
				"sing-box SOCKS port %s did not become ready: %w",
				address,
				waitCtx.Err(),
			)

		case <-ticker.C:
			conn, err := net.DialTimeout(
				"tcp",
				address,
				200*time.Millisecond,
			)

			if err == nil {
				_ = conn.Close()

				return nil
			}

			select {
			case processErr := <-p.done:
				if processErr != nil {
					return fmt.Errorf(
						"sing-box exited before becoming ready: %w",
						processErr,
					)
				}

				return fmt.Errorf(
					"sing-box exited before SOCKS port became ready",
				)

			default:
			}
		}
	}
}

func (p *Process) CheckIP(
	ctx context.Context,
	ipURL string,
	timeout time.Duration,
) (*CheckResult, error) {
	if p == nil {
		return nil, fmt.Errorf(
			"process is nil",
		)
	}

	if ipURL == "" {
		return nil, fmt.Errorf(
			"IP check URL is empty",
		)
	}

	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	proxyAddress := fmt.Sprintf(
		"127.0.0.1:%d",
		p.socksPort,
	)

	dialer, err := proxy.SOCKS5(
		"tcp",
		proxyAddress,
		nil,
		proxy.Direct,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create SOCKS5 dialer: %w",
			err,
		)
	}

	transport := &http.Transport{
		DialContext: func(
			dialCtx context.Context,
			network string,
			address string,
		) (net.Conn, error) {
			type dialResult struct {
				conn net.Conn
				err  error
			}

			resultCh := make(chan dialResult, 1)

			go func() {
				conn, err := dialer.Dial(
					network,
					address,
				)

				resultCh <- dialResult{
					conn: conn,
					err:  err,
				}
			}()

			select {
			case <-dialCtx.Done():
				return nil, dialCtx.Err()

			case result := <-resultCh:
				return result.conn, result.err
			}
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	requestCtx, cancel := context.WithTimeout(
		ctx,
		timeout,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(
		requestCtx,
		http.MethodGet,
		ipURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create IP request: %w",
			err,
		)
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return &CheckResult{
			Success: false,
			Latency: time.Since(start),
			Error:   err,
		}, nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(
		io.LimitReader(resp.Body, 1024),
	)
	if err != nil {
		return &CheckResult{
			Success: false,
			Latency: time.Since(start),
			Error:   err,
		}, nil
	}

	latency := time.Since(start)

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {
		return &CheckResult{
			Success: false,
			Latency: latency,
			Error: fmt.Errorf(
				"IP service returned HTTP %d",
				resp.StatusCode,
			),
		}, nil
	}

	return &CheckResult{
		Success: true,
		IP:      string(body),
		Latency: latency,
	}, nil
}

func (p *Process) Stop(
	timeout time.Duration,
) error {
	if p == nil {
		return nil
	}

	if p.cmd == nil ||
		p.cmd.Process == nil {
		_ = os.RemoveAll(p.tempDir)
		return nil
	}

	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	_ = p.cmd.Process.Signal(
		os.Interrupt,
	)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-p.done:
		_ = os.RemoveAll(p.tempDir)

		return nil

	case <-timer.C:
		_ = p.cmd.Process.Kill()

		select {
		case <-p.done:
		case <-time.After(time.Second):
		}

		_ = os.RemoveAll(p.tempDir)

		return nil
	}
}

func (p *Process) ConfigPath() string {
	if p == nil {
		return ""
	}

	return p.configPath
}

func (p *Process) SOCKSPort() int {
	if p == nil {
		return 0
	}

	return p.socksPort
}

func (p *Process) Logs() (string, string) {
	if p == nil {
		return "", ""
	}

	stdoutPath := filepath.Join(
		p.tempDir,
		"stdout.log",
	)

	stderrPath := filepath.Join(
		p.tempDir,
		"stderr.log",
	)

	stdout, _ := os.ReadFile(stdoutPath)
	stderr, _ := os.ReadFile(stderrPath)

	return string(stdout), string(stderr)
}
