# VPN Checker Backend

Go-бэкенд для сбора, проверки и хранения рабочих VPN-конфигураций.

Проект собирает конфигурации из нескольких источников, нормализует и удаляет дубликаты, выполняет быстрый TCP pre-check, проверяет конфигурации через `sing-box`, определяет внешний IP через VPN и сохраняет рабочие конфигурации для последующего использования и экспорта.

## Возможности

- Сбор конфигураций из нескольких источников.
- Поддержка VLESS, VMess, Trojan и Hysteria2.
- Нормализация и дедупликация конфигураций.
- TCP pre-check для быстрого отсечения недоступных endpoint.
- Проверка через `sing-box`.
- Реальный HTTP-запрос через VPN для подтверждения работоспособности.
- Определение внешнего IP.
- Измерение времени ответа через VPN.
- Параллельная проверка конфигураций.
- Realtime-прогресс, скорость и ETA.
- Web UI.
- REST API.
- Статистика по протоколам.
- Список лучших рабочих конфигураций.
- Экспорт TXT, JSON и Sing-box.
- Ручной запуск проверки через Web UI или API.
- Опциональный scheduler.

## Требования

- Go 1.22+.
- `sing-box`.
- Linux, macOS или Windows.
- Доступ в интернет.

Путь к `sing-box` задаётся в `configs/config.yaml`. Если бинарник доступен через `PATH`, достаточно указать `sing-box`.

Проверка установки:

```bash
go version
sing-box version
```

## Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
```

Скачайте зависимости:

```bash
go mod download
```

Проверьте проект:

```bash
go test ./...
```

## Запуск на Linux

### Arch Linux

```bash
sudo pacman -S go
```

Установите `sing-box` удобным для вашей системы способом и проверьте:

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

Откройте Web UI:

```text
http://localhost:8080/
```

Остановить сервер можно через `Ctrl+C`.

### Ubuntu / Debian

Установите Go из пакетов системы или с официального сайта Go:

```bash
sudo apt update
sudo apt install golang
```

После установки:

```bash
go version
sing-box version
```

Дальше используйте те же команды `go mod download`, `go run ./cmd/collector` и `go run ./cmd/server`.

## Запуск на macOS

При наличии Homebrew:

```bash
brew install go
```

Установите `sing-box` удобным способом и проверьте:

```bash
go version
sing-box version
```

Клонируйте проект и установите зависимости:

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

Web UI:

```text
http://localhost:8080/
```

Остановить сервер можно через `Ctrl+C`.

## Запуск на Windows

Установите Go с официального сайта:

```text
https://go.dev/dl/
```

Установите `sing-box` и добавьте его каталог в `PATH`, либо укажите полный путь в `configs/config.yaml`.

Проверка:

```powershell
go version
sing-box version
```

Клонирование и установка зависимостей:

```powershell
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
go mod download
```

Сбор конфигураций:

```powershell
go run .\cmd\collector
```

Запуск сервера:

```powershell
go run .\cmd\server
```

Web UI:

```text
http://localhost:8080/
```

Остановить сервер можно через `Ctrl+C`.

## Конфигурация

Основной файл:

```text
configs/config.yaml
```

Источники:

```text
configs/sources.yaml
```

В текущем наборе используются четыре источника `igareck` и FAST-наборы Russia и Europe из `kort0881/vpn-checker-backend`.

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

Linux/macOS с абсолютным путём:

```yaml
binary: /usr/local/bin/sing-box
```

Windows:

```yaml
binary: C:\sing-box\sing-box.exe
```

## Ручная проверка

Проверку можно запустить из Web UI кнопкой **«Запустить проверку»**.

Также доступен API:

```bash
curl -X POST http://localhost:8080/api/check
```

Текущий статус:

```bash
curl http://localhost:8080/api/check/status
```

Пример ответа во время проверки:

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

Web UI опрашивает этот endpoint в реальном времени и обновляет карточки, progress bar, скорость и ETA без перезагрузки страницы.

## Как определяется рабочая конфигурация

Для VLESS, VMess и Trojan сначала выполняется TCP pre-check endpoint:

```text
server:port
```

Если endpoint недоступен, конфигурация отбрасывается без запуска `sing-box`.

Для полноценной проверки запускается `sing-box` с временной конфигурацией. После готовности локального SOCKS5-прокси выполняется HTTP GET через VPN к сервису определения внешнего IP.

Конфигурация считается рабочей, если запрос через VPN успешно завершён и внешний IP получен.

Для Hysteria2 TCP pre-check не используется, поскольку проверка выполняется непосредственно через `sing-box`.

## Latency

Поле `latency` в результатах означает **время ответа через VPN**.

Это не ICMP ping и не чистый RTT до VPN-сервера. Измеряется время HTTP-запроса через локальный SOCKS5-прокси и VPN-туннель до сервиса определения внешнего IP.

В Web UI эта метрика подписана как:

```text
Время ответа через VPN
```

Она используется для сортировки раздела **«Лучшие конфигурации»**.

## Scheduler

Автоматический scheduler является опциональным.

Чтобы отключить автоматические проверки и запускать checker только вручную, укажите:

```yaml
health_check_interval: "0s"
```

При таком значении scheduler не запускается.

Проверку можно запускать только:

- через кнопку **«Запустить проверку»**;
- через `POST /api/check`.

Это также означает, что обновление Web UI не запускает новую проверку.

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

### Web UI

```text
http://localhost:8080/
```

Web UI включает:

- текущую статистику проверки;
- количество рабочих и нерабочих конфигураций;
- realtime progress bar;
- скорость проверки;
- ETA;
- статистику по протоколам;
- лучшие конфигурации;
- экспорт результатов;
- Catppuccin Latte для светлой темы;
- Catppuccin Mocha для тёмной темы.

Тема переключается автоматически в соответствии с системной настройкой `prefers-color-scheme`.

## Экспорт

Сохранённые рабочие конфигурации находятся в:

```text
data/all-working.txt
```

Скачать TXT:

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

В Web UI экспорт доступен отдельными кнопками:

```text
Скачать TXT
Скачать JSON
Скачать Sing-box
```

Экспорт не запускает новую проверку — используются уже сохранённые результаты.

## Результаты и данные

```text
data/configs/all.txt   # актуальный входной набор

data/results.json      # результаты проверок

data/stats.json        # сохранённая статистика

data/all-working.txt   # рабочие конфигурации
```

Файлы в `data/` могут содержать результаты конкретного запуска и не должны рассматриваться как вечный кэш доступности VPN.

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
```

## Тесты и форматирование

Запустить все тесты:

```bash
go test ./...
```

Форматирование Go-кода:

```bash
gofmt -w .
```

## Пример контрольного запуска

После перехода от большого сырого набора к FAST-источникам был получен контрольный результат:

```text
2047 конфигураций
84 рабочих
1963 нерабочих
9 мин 16 сек
3.68 cfg/s
```

Главное ускорение достигнуто за счёт более качественного входного набора и TCP pre-check, а не за счёт увеличения количества повторных сетевых запросов.

## Лицензия

Проект распространяется по лицензии **MIT License**.

MIT разрешает использовать, копировать, изменять, объединять, публиковать, распространять, сублицензировать и продавать копии проекта при сохранении уведомления об авторских правах и текста лицензии.

Полный текст лицензии находится в файле [`LICENSE`](LICENSE).

## Примечания

VPN-конфигурации являются динамическими: рабочая конфигурация сегодня может перестать работать позже.

Поэтому рабочие результаты нужно периодически перепроверять на актуальном наборе источников.
