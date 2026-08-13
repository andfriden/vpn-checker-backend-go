# VPN Checker Backend

[🇬🇧 English](README.md) | [🇷🇺 Русский](README.ru.md)

Go-бэкенд для сбора, проверки, хранения и экспорта рабочих VPN-конфигураций.

Проект собирает VPN-конфигурации из нескольких источников, нормализует и удаляет дубликаты, выполняет быстрый TCP pre-check, запускает `sing-box` для полноценной проверки VPN, проверяет реальный выход в интернет через туннель, измеряет время ответа и сохраняет рабочие конфигурации.

## Возможности

* Поддержка VLESS, VMess, Trojan и Hysteria2.
* Сбор конфигураций из нескольких источников.
* Нормализация и дедупликация.
* TCP pre-check перед запуском `sing-box`.
* Полноценная проверка через `sing-box`.
* Реальный HTTP-запрос через VPN.
* Определение внешнего IP.
* Измерение времени ответа через VPN.
* Параллельная проверка конфигураций.
* Realtime-прогресс, скорость и ETA.
* REST API.
* Web UI.
* Catppuccin Latte и Mocha.
* Статистика по протоколам.
* Рейтинг лучших рабочих конфигураций.
* Экспорт TXT, JSON и Sing-box.
* Ручной запуск проверки через Web UI или API.
* Опциональный scheduler.
* Сборка релизов для Linux, macOS и Windows.
* SLSA provenance для релизных артефактов.

## Архитектура

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

## Требования

* Go 1.22+
* `sing-box`
* Linux, macOS или Windows
* Доступ в интернет

Проверить установку:

```bash
go version
sing-box version
```

## Установка

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Проверить проект:

```bash
go test ./...
```

## Запуск в Linux

### Arch Linux

```bash
sudo pacman -S go
```

Установите `sing-box` удобным для вашей системы способом.

Проверьте:

```bash
go version
sing-box version
```

Соберите актуальный набор конфигураций:

```bash
go run ./cmd/collector
```

Запустите сервер:

```bash
go run ./cmd/server
```

Web UI:

```text
http://localhost:8080/
```

Остановить сервер можно через `Ctrl+C`.

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang
```

Затем:

```bash
go version
sing-box version
go mod download
go run ./cmd/collector
go run ./cmd/server
```

## Запуск в macOS

При наличии Homebrew:

```bash
brew install go
```

Установите `sing-box` и проверьте:

```bash
go version
sing-box version
```

Клонируйте проект:

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Соберите конфигурации:

```bash
go run ./cmd/collector
```

Запустите сервер:

```bash
go run ./cmd/server
```

Откройте:

```text
http://localhost:8080/
```

## Запуск в Windows

Установите Go:

```text
https://go.dev/dl/
```

Установите `sing-box` и добавьте его в `PATH`, либо укажите полный путь в `configs/config.yaml`.

Проверьте:

```powershell
go version
sing-box version
```

Клонируйте проект:

```powershell
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Соберите конфигурации:

```powershell
go run .\cmd\collector
```

Запустите сервер:

```powershell
go run .\cmd\server
```

Откройте:

```text
http://localhost:8080/
```

Остановить сервер можно через `Ctrl+C`.

## Конфигурация

Основной конфигурационный файл:

```text
configs/config.yaml
```

Источники:

```text
configs/sources.yaml
```

После изменения источников заново соберите входной набор:

```bash
go run ./cmd/collector
```

Результат сохраняется в:

```text
data/configs/all.txt
```

### Путь к sing-box

Если `sing-box` находится в `PATH`:

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

## Ручная проверка

Проверка запускается в Web UI кнопкой **«Запустить проверку»**.

Также доступен API:

```bash
curl -X POST http://localhost:8080/api/check
```

Статус текущей проверки:

```bash
curl http://localhost:8080/api/check/status
```

Пример:

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

Web UI обновляет карточки, progress bar, скорость и ETA автоматически и не требует перезагрузки страницы.

## Как определяется рабочая конфигурация

Для VLESS, VMess и Trojan сначала выполняется TCP pre-check:

```text
server:port
```

Если endpoint недоступен, конфигурация отбрасывается без запуска `sing-box`.

После этого запускается `sing-box` с временной конфигурацией.

Когда локальный SOCKS5 становится доступен, выполняется HTTP GET через VPN к сервису определения внешнего IP.

Конфигурация считается рабочей, если:

* `sing-box` запустился;
* SOCKS5 стал доступен;
* HTTP-запрос через VPN успешно завершился;
* был получен внешний IP.

Для Hysteria2 TCP pre-check не используется, проверка выполняется непосредственно через `sing-box`.

## Latency

Поле `latency` означает **время ответа через VPN**.

Это не ICMP ping и не чистый RTT до VPN-сервера.

Измеряется:

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

В Web UI эта метрика обозначена как:

```text
Время ответа через VPN
```

Она используется для сортировки раздела **«Лучшие конфигурации»**.

## Scheduler

Автоматический scheduler можно отключить:

```yaml
health_check_interval: "0s"
```

При этом проверка запускается только вручную через Web UI или API:

```bash
curl -X POST http://localhost:8080/api/check
```

Обновление страницы не запускает новую проверку.

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

Frontend расположен в:

```text
web/
├── index.html
├── app.js
└── style.css
```

Web UI показывает:

* текущий статус проверки;
* количество рабочих и нерабочих конфигураций;
* progress bar;
* скорость;
* ETA;
* статистику по протоколам;
* лучшие конфигурации;
* экспорт;
* Catppuccin Latte в светлой теме;
* Catppuccin Mocha в тёмной теме.

Тема автоматически выбирается по системной настройке `prefers-color-scheme`.

## Экспорт

Сохранённые рабочие конфигурации находятся в:

```text
data/all-working.txt
```

Поддерживаются:

```text
TXT
JSON
Sing-box
```

API:

```text
GET /api/export
GET /api/export/json
GET /api/export/singbox
```

TXT:

```bash
curl -OJ http://localhost:8080/api/export
```

JSON:

```bash
curl -o results.json http://localhost:8080/api/export/json
```

Sing-box:

```bash
curl -o singbox.json http://localhost:8080/api/export/singbox
```

Экспорт использует уже сохранённые результаты и не запускает новую проверку.

## Результаты

```text
data/configs/all.txt   # актуальный входной набор
data/results.json      # результаты проверок
data/stats.json        # статистика
data/all-working.txt   # рабочие конфигурации
```

Результаты конкретного запуска не следует рассматривать как постоянный кэш доступности VPN.

## Релизы

GitHub Actions собирает релизные артефакты для:

```text
Linux amd64
Linux arm64

macOS amd64
macOS arm64

Windows amd64
```

Релиз запускается созданием Git-тега:

```bash
git tag <version>
git push origin <version>
```

Workflow собирает сервер и Collector и публикует архивы:

```text
vpn-checker-linux-amd64.tar.gz
vpn-checker-linux-arm64.tar.gz
vpn-checker-darwin-amd64.tar.gz
vpn-checker-darwin-arm64.tar.gz
vpn-checker-windows-amd64.zip
```

`sing-box` распространяется отдельно и должен быть установлен на целевой системе.

## SLSA Provenance

Для релизных артефактов используется SLSA Generic Generator.

Workflow:

1. скачивает реальные артефакты Release;
2. вычисляет SHA-256;
3. передаёт hashes в SLSA Generator;
4. создаёт provenance;
5. добавляет provenance в GitHub Release.

Это позволяет проверять происхождение опубликованных артефактов.

## Структура проекта

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

## Тесты

```bash
go test ./...
```

Форматирование:

```bash
gofmt -w .
```

## Производительность

Контрольный запуск FAST-набора:

```text
2047 конфигураций
84 рабочих
1963 нерабочих
9 мин 16 сек
3.68 cfg/s
```

Основное ускорение достигнуто за счёт более качественного входного набора, дедупликации и TCP pre-check.

## Лицензия

Проект распространяется под **MIT License**.

MIT разрешает использовать, копировать, изменять, объединять, публиковать, распространять, выдавать сублицензии и продавать копии программного обеспечения при сохранении уведомления об авторских правах и текста лицензии.

Полный текст находится в файле [`LICENSE`](LICENSE).

