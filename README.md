# VPN Checker Backend

Go-бэкенд для сбора, проверки и хранения рабочих VPN-конфигураций.

Проект загружает конфигурации из нескольких источников, удаляет дубли, выполняет быстрый TCP pre-check, запускает `sing-box` для полноценной проверки и сохраняет только рабочие конфигурации.

## Возможности

* Сбор VPN-конфигураций из нескольких источников.
* Поддержка VLESS, VMess, Trojan и Hysteria2.
* Удаление дубликатов перед проверкой.
* TCP pre-check для быстрого отсечения недоступных endpoint.
* Проверка конфигурации через `sing-box`.
* Проверка реального выхода в интернет через VPN.
* Определение внешнего IP.
* Измерение времени ответа через VPN.
* Параллельная проверка конфигураций.
* Realtime-прогресс, скорость и ETA.
* JSON API.
* Web-интерфейс.
* Экспорт рабочих конфигураций.
* Статистика по протоколам.
* Список лучших конфигураций.

## Архитектура

```text
sources.yaml
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
          Web / API
```

## Требования

* Go 1.22+
* `sing-box`
* macOS или Linux
* доступ в интернет

Путь к `sing-box` указывается в `configs/config.yaml`.

Проверить установку:

```bash
sing-box version
```

## Установка

Клонировать репозиторий:

```bash
git clone https://github.com/andfriden/vpn-checker-backend-go.git
cd vpn-checker-backend-go
```

Скачать зависимости:

```bash
go mod download
```

Проверить проект:

```bash
go test ./...
```

## Конфигурация

Основной конфигурационный файл:

```text
configs/config.yaml
```

Источники:

```text
configs/sources.yaml
```



Собрать актуальный набор конфигураций:

```bash
go run ./cmd/collector
```

Результат будет сохранён в:

```text
data/configs/all.txt
```

## Запуск сервера

```bash
go run ./cmd/server
```

После запуска Web UI доступен по адресу:

```text
http://localhost:8080/
```

API:

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

## Проверка конфигураций

Проверка запускается из Web UI кнопкой:

```text
Запустить проверку
```

или через API:

```bash
curl -X POST http://localhost:8080/api/check
```

Текущий статус:

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

Web-интерфейс обновляет эти значения автоматически без перезагрузки страницы.

## Как определяется рабочий конфиг

Для большинства протоколов сначала выполняется TCP pre-check:

```text
server:port
```

Если endpoint недоступен, конфигурация сразу считается нерабочей.

После этого запускается `sing-box` с временной конфигурацией.

Когда SOCKS5 становится доступен, выполняется HTTP-запрос через VPN к сервису определения IP.

Конфигурация считается рабочей, если запрос успешно завершился и был получен внешний IP.

## Latency

Значение `latency` в результатах — это время выполнения HTTP-запроса через VPN.

То есть измеряется:

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

Это не ICMP ping и не чистый RTT до VPN-сервера.

В Web UI эта метрика отображается как:

```text
Время ответа через VPN
```

## Результаты

Рабочие конфигурации:

```text
data/all-working.txt
```

Все результаты проверки:

```text
data/results.json
```

Статистика:

```text
data/stats.json
```

## Экспорт

Скачать сохранённые рабочие конфигурации:

```bash
curl -OJ http://localhost:8080/api/export
```

Web UI также предоставляет экспорт сохранённого `all-working.txt`.

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

Frontend находится в:

```text
web/
├── index.html
├── app.js
└── style.css
```

Интерфейс поддерживает:

* Catppuccin Latte для светлой темы.
* Catppuccin Mocha для тёмной темы.
* автоматическое переключение по системной теме.
* realtime progress bar.
* текущую скорость проверки.
* ETA.
* количество рабочих и нерабочих конфигураций.
* статистику по протоколам.
* лучшие конфигурации.
* экспорт результатов.

## Scheduler

Автоматический scheduler можно отключить:

```yaml
health_check_interval: "0s"
```

В этом режиме проверка запускается только вручную через Web UI или API.

Это удобно для VPN-конфигураций, поскольку рабочее состояние конфигураций может быстро меняться и результаты старого запуска не следует считать актуальными бесконечно.

## Производительность

После перехода с большого сырого набора примерно в 34 000 конфигураций на FAST-набор из примерно 2 000 конфигураций удалось существенно сократить время проверки.

Контрольный запуск:

```text
2047 конфигураций
84 рабочих
9 мин 16 сек
3.68 cfg/s
```

Основное ускорение достигается не увеличением количества worker'ов, а использованием более качественных источников и TCP pre-check.

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

VPN-конфигурации являются динамическими: рабочий сегодня конфиг может перестать работать позже.

Поэтому результаты проверки не являются постоянным кэшем доступности. Каждый новый запуск должен проверять актуальный набор конфигураций заново.

## License

This project is licensed under the MIT License.
