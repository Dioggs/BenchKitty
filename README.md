<div align="center">

# 🐱 BenchKitty

**A dead-simple HTTP benchmarking tool for cat enthusiasts.**

Give it a URL and watch the magic happen.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey)
![Made with love](https://img.shields.io/badge/made%20with-%E2%9D%A4%20and%20cats-ff69b4)

</div>

---

## What is BenchKitty?

BenchKitty fires a configurable number of HTTP requests at an endpoint, measures
how long each one takes, and turns the raw latency data into a report
either in your terminal or as an Excel workbook.

It is intentionally small: no config files, no dependencies at runtime, just a
URL and a handful of flags.

## Why was BenchKitty made?

BenchKitty was made with the intention of being a very simple project to practice some Golang concepts, so don't treat it like the Linux kernel

## Features

- Any common HTTP method — `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- Request bodies for `POST` / `PUT` / `PATCH` (sent as `application/json`)
- Per-request latency plus `Avg`, `P50`, `P95` and `P99`
- Throughput and requests-per-second metrics
- Status code breakdown of every response
- Optional Excel report with summary, latency and status analyses + charts
- Optional pretty-print of each individual call

## Installation

Requires [Go](https://go.dev/) 1.26+.

```bash
git clone <your-repo-url> BenchKitty
cd BenchKitty
go build -o BenchKitty .
```

Then run it:

```bash
./BenchKitty -r 100 -d 500 https://example.com
```

Or, without building:

```bash
go run . -r 100 -d 500 https://example.com
```

## Usage

```bash
BenchKitty [flags] <url>
```

### Flags

| Flag | Description                                                              | Default    |
| :--- | :----------------------------------------------------------------------- | :--------- |
| `-r` | Number of requests to send                                               | `100`      |
| `-d` | Delay between every request, in milliseconds                             | `1000`     |
| `-t` | HTTP method: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`                     | `GET`      |
| `-b` | Request body — only valid for `POST`, `PUT`, `PATCH` (sent as JSON)      | *(none)*   |
| `-o` | Output **directory** for `benchmark.xlsx` (defaults to the terminal)     | *(none)*   |
| `-p` | Pretty-print every individual job                                        | `false`    |
| `<url>` | Target endpoint (required)                                            | —          |

> Passing `-b` with a method that cannot carry a body (like `GET`) is
> rejected before the benchmark starts.

### Examples

```bash
# Simple GET benchmark
BenchKitty https://example.com

# 200 POSTs with a JSON body, 100ms apart
BenchKitty -r 200 -d 100 -t POST \
  -b '{"title":"foo","body":"bar","userId":1}' \
  https://example.com/posts

# Hammer an endpoint with no delay and save an Excel report
BenchKitty -r 500 -d 0 -o ./reports https://example.com

# Watch every single call as it happens
BenchKitty -p -r 10 -t GET https://example.com
```

## The report

After the run, BenchKitty prints a summary like this:

```text
{Avg:26 P50:21 P95:50 P99:50 RPS:18 TPS:1517 StatusCount:map[200:5]}
```

| Metric        | Meaning                                                          |
| :------------ | :--------------------------------------------------------------- |
| `Avg`         | Average latency (ms)                                             |
| `P50`         | Median latency (ms)                                              |
| `P95` / `P99` | Latency at the 95th / 99th percentile (ms)                       |
| `RPS`         | Requests per second                                              |
| `TPS`         | Throughput in bytes per second                                   |
| `StatusCount` | How many responses came back with each HTTP status code          |

### Excel output

When `-o <dir>` is given, BenchKitty writes `<dir>/benchmark.xlsx` with three
sheets:

| Sheet      | Contents                                                        |
| :--------- | :------------------------------------------------------------- |
| `Summary`  | All headline metrics plus a percentile chart                    |
| `Latency`  | Per-request latency, sorted distribution and line/bar charts    |
| `Status`   | Status code distribution with a pie chart                       |

> The file always uses the same name, so re-running with the same `-o`
> **overwrites** the previous `benchmark.xlsx`.

## 🛠️ Development

Common tasks are wrapped in the `Makefile`:

```bash
make run     # quick smoke run against a sample endpoint
make build   # compile the BenchKitty binary
make test    # run the test suite
make fmt     # format the code with gofmt
make vet     # run static checks with go vet
make check   # fmt + vet + test
make clean   # remove build artifacts
```

Run the test suite:

```bash
go test ./...
```

Useful extras for contributors:

```bash
gofmt -l .     # formatting
go vet ./...   # static checks
go build ./... # compile everything
```

