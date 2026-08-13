# VPN Checker Backend

[🇬🇧 English](README.md) | [🇷🇺 Русский](README.ru.md)

A Go backend for collecting, checking, storing, and exporting working VPN configurations.

The project collects VPN configurations from multiple sources, normalizes and deduplicates them, performs a fast TCP pre-check, launches `sing-box` for a full VPN test, verifies real Internet access through the tunnel, measures response time, and stores working configurations.

## Features

* VLESS, VMess, Trojan, and Hysteria2 support.
* Collect configurations from multiple sources.
* Normalize and deduplicate configurations.
* Fast TCP pre-check before running `sing-box`.
* Full VPN validation through `sing-box`.
* Real HTTP request through the VPN tunnel.
* External IP detection.
* VPN response-time measurement.
* Parallel configuration checks.
* Realtime progress, speed, and ETA.
* REST API.
* Web UI.
* Catppuccin Latte and Mocha themes.
* Protocol statistics.
* Best working configuration ranking.
* TXT, JSON, and Sing-box export.
* Manual checks through Web UI or API.
* Optional scheduler.
* Cross-platform release builds for Linux, macOS, and Windows.
* SLSA provenance for release artifacts.

## Architecture

```text
configs/sources.yaml
        │
        ▼
    Collector
        │
        ▼
 Normalize / Deduplicate
        │
        ▼
 data/configs/all.txt
        │
        ▼
    TCP pre-check
        │
        ▼
      sing-box
        │
        ▼
 HTTP request through VPN
        │
        ▼
    CheckResult
        │
        ├── data/results.json
        ├── data/stats.json
        └── data/all-working.txt
                │
                ▼
             API / Web UI
```

## Requirements

* Go 1.22+
* `sing-box`
* Linux, macOS, or Windows
* Internet access

Check the installation:

```bash
go version
sing-box version
```

## Installation

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Run tests:

```bash
go test ./...
```

## Running on Linux

### Arch Linux

```bash
sudo pacman -S go
```

Install `sing-box` using the method appropriate for your system.

Build the current configuration set:

```bash
go run ./cmd/collector
```

Start the server:

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080/
```

Stop the server with `Ctrl+C`.

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang
```

Then:

```bash
go version
sing-box version
go mod download
go run ./cmd/collector
go run ./cmd/server
```

## Running on macOS

With Homebrew:

```bash
brew install go
```

Install `sing-box`, then verify:

```bash
go version
sing-box version
```

Clone and prepare the project:

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Build the configuration set:

```bash
go run ./cmd/collector
```

Start the server:

```bash
go run ./cmd/server
```

Open:

```text
http://localhost:8080/
```

## Running on Windows

Install Go from:

```text
https://go.dev/dl/
```

Install `sing-box` and add it to `PATH`, or specify its full path in `configs/config.yaml`.

Check:

```powershell
go version
sing-box version
```

Clone the project:

```powershell
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Build the configuration set:

```powershell
go run .\cmd\collector
```

Start the server:

```powershell
go run .\cmd\server
```

Open:

```text
http://localhost:8080/
```

## Configuration

Main configuration:

```text
configs/config.yaml
```

Source configuration:

```text
configs/sources.yaml
```

After changing the sources, rebuild the input set:

```bash
go run ./cmd/collector
```

The result is stored in:

```text
data/configs/all.txt
```

### sing-box path

If `sing-box` is available in `PATH`:

```yaml
binary: sing-box
```

Linux/macOS:

```yaml
binary: /usr/local/bin/sing-box
```

Windows:

```yaml
binary: C:\sing-box\sing-box.exe
```

## Running a Check

Start a check from the Web UI using **Start Check**, or through the API:

```bash
curl -X POST http://localhost:8080/api/check
```

Get the current status:

```bash
curl http://localhost:8080/api/check/status
```

Example:

```json
{
  "status": "running",
  "total": 2047,
  "checked": 1161,
  "working": 67,
  "failed": 1094,
  "progress": 56,
  "current_speed": 3.58,
  "estimated_seconds_left": 247.31
}
```

The Web UI updates the progress bar, counters, speed, and ETA automatically without reloading the page.

## How a Configuration Is Validated

For VLESS, VMess, and Trojan, the checker first performs a TCP pre-check:

```text
server:port
```

If the endpoint is unreachable, the configuration is rejected before starting `sing-box`.

For the full check, `sing-box` is started with a temporary configuration. Once the local SOCKS5 proxy becomes ready, the checker performs an HTTP request through the VPN tunnel to an external IP service.

A configuration is considered working when:

* `sing-box` starts successfully;
* the local SOCKS5 proxy becomes available;
* the HTTP request through the VPN succeeds;
* an external IP is returned.

Hysteria2 skips the TCP pre-check and is tested directly through `sing-box`.

## Latency

The `latency` value is the **response time through the VPN**.

It is not an ICMP ping and not a pure RTT to the VPN server.

Measured path:

```text
SOCKS5
   ↓
VPN tunnel
   ↓
HTTP request
   ↓
External IP service
   ↓
Response
```

The Web UI labels this metric:

```text
VPN response time
```

It is used to rank the **Best Configurations** section.

## Results

Working configurations:

```text
data/all-working.txt
```

All check results:

```text
data/results.json
```

Statistics:

```text
data/stats.json
```

## Export

The Web UI supports:

```text
Download TXT
Download JSON
Download Sing-box
```

API endpoints:

```text
GET /api/export
GET /api/export/json
GET /api/export/singbox
```

Download TXT:

```bash
curl -OJ http://localhost:8080/api/export
```

Download JSON:

```bash
curl -o results.json http://localhost:8080/api/export/json
```

Download Sing-box configuration:

```bash
curl -o singbox.json http://localhost:8080/api/export/singbox
```

Export uses already stored results and does not start a new check.

## Scheduler

Automatic checks can be disabled with:

```yaml
health_check_interval: "0s"
```

In this mode, checks are started only manually through the Web UI or API.

```bash
curl -X POST http://localhost:8080/api/check
```

Reloading the Web UI does not start a new check.

## API

```text
GET  /health

GET  /api/results
GET  /api/best
GET  /api/best/singbox

POST /api/check
GET  /api/check/status

GET  /api/stats

GET  /api/export
GET  /api/export/json
GET  /api/export/singbox
```

## Web UI

Frontend files:

```text
web/
├── index.html
├── app.js
└── style.css
```

The Web UI provides:

* realtime check status;
* working/failed counters;
* progress bar;
* current speed;
* ETA;
* protocol statistics;
* best working configurations;
* TXT, JSON, and Sing-box export;
* Catppuccin Latte for light mode;
* Catppuccin Mocha for dark mode.

The theme follows the system `prefers-color-scheme` setting.

## Releases

GitHub Actions can build release artifacts for:

```text
Linux amd64
Linux arm64

macOS amd64
macOS arm64

Windows amd64
```

Create a release by pushing a version tag:

```bash
git tag <version>
git push origin <version>
```

The release workflow creates archives containing:

* the VPN Checker server;
* the Collector binary.

Example artifact names:

```text
vpn-checker-linux-amd64.tar.gz
vpn-checker-linux-arm64.tar.gz
vpn-checker-darwin-amd64.tar.gz
vpn-checker-darwin-arm64.tar.gz
vpn-checker-windows-amd64.zip
```

`sing-box` is distributed separately and must be installed on the target system.

## SLSA Provenance

Release artifacts are processed by the SLSA Generic Generator.

The workflow:

1. downloads the actual release artifacts;
2. calculates SHA-256 digests;
3. generates SLSA provenance;
4. uploads the provenance to the GitHub Release.

This allows users to verify the origin of published release artifacts.

## Project Structure

```text
cmd/
├── collector/
│   └── main.go
└── server/
    └── main.go

configs/
├── config.yaml
└── sources.yaml

data/
├── all-working.txt
├── configs/
│   └── all.txt
├── results.json
└── stats.json

internal/
├── api/
├── app/
├── checker/
├── config/
├── decoder/
├── downloader/
├── extractor/
├── model/
├── normalizer/
├── parser/
├── pipeline/
├── precheck/
├── scheduler/
├── service/
├── singbox/
├── source/
├── storage/
├── validator/
└── worker/

web/
├── app.js
├── index.html
└── style.css

.github/
└── workflows/
    ├── release.yml
    └── generator-generic-ossf-slsa3-publish.yml
```

## Tests

Run all tests:

```bash
go test ./...
```

Format Go code:

```bash
gofmt -w .
```

## License

This project is licensed under the MIT License.

See [`LICENSE`](LICENSE) for the full license text.

MIT permits using, copying, modifying, merging, publishing, distributing, sublicensing, and selling copies of the software, provided that the copyright notice and license text are preserved.

