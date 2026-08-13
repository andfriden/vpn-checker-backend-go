# VPN Checker Backend

Go-бэкенд для сбора, проверки и хранения рабочих VPN-конфигураций.

Проект собирает конфигурации из нескольких источников, удаляет дубли, выполняет быстрый TCP pre-check, запускает `sing-box` для полноценной проверки VPN и сохраняет рабочие конфигурации.

## Возможности

* Поддержка VLESS, VMess, Trojan и Hysteria2.
* Сбор конфигураций из нескольких источников.
* Нормализация и удаление дубликатов.
* TCP pre-check перед запуском `sing-box`.
* Полноценная проверка через `sing-box`.
* Проверка реального выхода в интернет через VPN.
* Определение внешнего IP.
* Измерение времени ответа через VPN.
* Параллельная проверка конфигураций.
* Realtime-прогресс, скорость и ETA.
* JSON API.
* Web UI.
* Catppuccin Latte / Mocha.
* Статистика по протоколам.
* Рейтинг лучших рабочих конфигураций.
* Экспорт TXT, JSON и Sing-box.
* Автоматическая сборка релизов для Linux, macOS и Windows.
* SLSA provenance для release-артефактов.

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

## Установка из исходников

Клонировать репозиторий:

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
```

Установить зависимости:

```bash
go mod download
```

Проверить проект:

```bash
go test ./...
```

## Запуск на Linux

### Arch Linux

```bash
sudo pacman -S go
```

Установить `sing-box` удобным для системы способом.

Проверить:

```bash
go version
sing-box version
```

Собрать список конфигураций:

```bash
go run ./cmd/collector
```

Запустить сервер:

```bash
go run ./cmd/server
```

Открыть:

```text
http://localhost:8080/
```

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang
```

Затем:

```bash
go mod download
go run ./cmd/collector
go run ./cmd/server
```

## Запуск на macOS

Установить Go через Homebrew:

```bash
brew install go
```

Проверить:

```bash
go version
sing-box version
```

Собрать список:

```bash
go run ./cmd/collector
```

Запустить:

```bash
go run ./cmd/server
```

Web UI:

```text
http://localhost:8080/
```

## Запуск на Windows

Установить Go с:

```text
https://go.dev/dl/
```

Проверить в PowerShell:

```powershell
go version
sing-box version
```

Клонировать репозиторий:

```powershell
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
```

Установить зависимости:

```powershell
go mod download
```

Собрать конфигурации:

```powershell
go run .\cmd\collector
```

Запустить сервер:

```powershell
go run .\cmd\server
```

Открыть:

```text
http://localhost:8080/
```

Остановить сервер:

```text
Ctrl+C
```

## Настройка sing-box

Путь к бинарнику задаётся в:

```text
configs/config.yaml
```

Если `sing-box` находится в `PATH`:

```yaml
binary: sing-box
```

Linux/macOS с абсолютным путём:

```yaml
binary: /usr/local/bin/sing-box
```

Windows:

```yaml
binary: C:\sing-box\sing-box.exe
```

## Источники конфигураций

Источники задаются в:

```text
configs/sources.yaml
```

Текущая конфигурация использует FAST-наборы Russia и Europe вместе с дополнительными источниками.

Сбор:

```bash
go run ./cmd/collector
```

Результат:

```text
data/configs/all.txt
```

## Проверка конфигураций

Проверку можно запустить через Web UI:

```text
Запустить проверку
```

или через API:

```bash
curl -X POST http://localhost:8080/api/check
```

Статус:

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
  "started_at": "2026-08-13T12:38:02+03:00",
  "elapsed_seconds": 324.07,
  "current_speed": 3.58,
  "estimated_seconds_left": 247.31
}
```

Web UI обновляет прогресс автоматически без перезагрузки страницы.

## Как определяется рабочий конфиг

Для большинства протоколов сначала выполняется TCP pre-check:

```text
server:port
```

Если endpoint недоступен, конфигурация отсеивается без запуска полного VPN-теста.

После этого запускается `sing-box` с временной конфигурацией.

Когда SOCKS5 становится доступен, выполняется HTTP-запрос через VPN к сервису определения внешнего IP.

Конфигурация считается рабочей, если:

* `sing-box` успешно запустился;
* SOCKS5 стал доступен;
* HTTP-запрос через VPN завершился успешно;
* получен внешний IP.

## Latency

`latency` — это время выполнения HTTP-запроса через VPN.

Это **не ICMP ping** и не чистый RTT до VPN-сервера.

Схема:

```text
SOCKS5
   ↓
VPN tunnel
   ↓
HTTP request
   ↓
external IP service
   ↓
response
```

В Web UI эта метрика обозначена как:

```text
Время ответа через VPN
```

Именно она используется для сортировки лучших рабочих конфигураций.

## Результаты

Рабочие конфигурации:

```text
data/all-working.txt
```

Все результаты:

```text
data/results.json
```

Статистика:

```text
data/stats.json
```

## Экспорт

Web UI поддерживает три варианта экспорта:

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

Пример TXT:

```bash
curl -OJ http://localhost:8080/api/export
```

JSON:

```bash
curl -o results.json \
  http://localhost:8080/api/export/json
```

Sing-box:

```bash
curl -o singbox.json \
  http://localhost:8080/api/export/singbox
```

## Web UI

Frontend:

```text
web/
├── index.html
├── app.js
└── style.css
```

Интерфейс поддерживает:

* Catppuccin Latte для светлой темы.
* Catppuccin Mocha для тёмной темы.
* автоматический выбор темы через `prefers-color-scheme`.
* realtime progress bar.
* текущую скорость проверки.
* ETA.
* количество рабочих и нерабочих конфигураций.
* статистику по протоколам.
* лучшие конфигурации.
* экспорт результатов.

## API

Основные endpoints:

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

## Scheduler

Автоматический scheduler можно отключить:

```yaml
health_check_interval: "0s"
```

При таком значении проверка не запускается автоматически.

Проверка выполняется только вручную:

```bash
curl -X POST http://localhost:8080/api/check
```

или через кнопку Web UI.

Это позволяет каждый раз проверять актуальный набор конфигураций вручную.

## Производительность

Переход от большого сырого набора примерно в 34 000 конфигураций к FAST-источникам примерно на 2 000 конфигураций существенно сократил время проверки.

Контрольный запуск:

```text
2047 конфигураций
84 рабочих
1963 нерабочих
4.10% рабочих
9 мин 16 сек
3.68 cfg/s
```

Основное ускорение обеспечивается:

* более качественными источниками;
* удалением дублей;
* TCP pre-check;
* параллельной обработкой.

## Релизы

Проект поддерживает автоматическую сборку release-артефактов через GitHub Actions.

Релиз запускается созданием Git-тега:

```bash
git tag v0.1.0
git push origin v0.1.0
```

После этого workflow собирает:

```text
Linux amd64
Linux arm64

macOS amd64
macOS arm64

Windows amd64
```

Для каждой платформы создаются архивы с:

* `vpn-checker` server;
* `vpn-checker-collector`.

Примеры названий:

```text
vpn-checker-linux-amd64.tar.gz
vpn-checker-linux-arm64.tar.gz
vpn-checker-darwin-amd64.tar.gz
vpn-checker-darwin-arm64.tar.gz
vpn-checker-windows-amd64.zip
```

Готовые бинарники публикуются в GitHub Releases.

### sing-box

`sing-box` не встраивается в бинарник VPN Checker и должен быть установлен отдельно.

## SLSA Provenance

Для release-артефактов используется SLSA Generic Generator.

После создания GitHub Release workflow:

1. скачивает реальные release-артефакты;
2. вычисляет SHA-256;
3. передаёт hashes в SLSA Generator;
4. создаёт provenance;
5. добавляет provenance в GitHub Release.

Это позволяет проверять происхождение опубликованных release-артефактов.

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

Запустить все тесты:

```bash
go test ./...
```

Форматирование:

```bash
gofmt -w .
```

## Примечания

VPN-конфигурации являются динамическими.

Рабочая конфигурация может перестать работать в любой момент, поэтому результаты проверки не следует считать постоянным кэшем доступности.

Каждый новый запуск рекомендуется выполнять на актуальном наборе конфигураций.

## License

This project is licensed under the MIT License.

Полный текст лицензии находится в файле:

```text
LICENSE
```

MIT License разрешает использовать, изменять, распространять и создавать производные версии проекта при сохранении уведомления об авторских правах и текста лицензии.

