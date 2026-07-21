# GitHub Webhook Listener

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

A lightweight Go service for receiving GitHub Webhooks and executing Shell commands. Includes an optional built-in Web dashboard with project health monitoring and webhook execution logs.

| <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00001.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00002.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00003.png"> |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |

## Features

- **Webhook Receiver**: Handles GitHub events such as `push`, `pull_request`, `release`, with HMAC-SHA256 signature verification
- **Rule Matching**: Flexible rules by event type and branch, triggering Shell commands with optional retry on failure
- **Web Dashboard**: Optional dashboard (with Basic Auth) for viewing project health status, 24-hour uptime charts, and webhook execution logs
- **Health Monitoring**: Per-repository scheduled HTTP probes with status tracking
- **Single Binary**: Compiled as a standalone Go executable with embedded SQLite — no external database required

## Quick Start

### 1. Get the binary

Download the binary for your platform from [Releases](https://github.com/zxc7563598/github-webhook-listener/releases), or build locally:

```bash
make build          # Build to bin/ (current platform)
make build-linux    # Linux amd64
make build-darwin   # macOS amd64 + arm64
make build-windows  # Windows amd64
```

> [!NOTE]
> Local builds require Go 1.22+. Cross-compilation or Web dashboard features also require Node.js (for building the frontend).

### 2. Create the config file

```bash
cp config/config.example.yaml config.yaml
```

Edit `config.yaml` with your repository and Webhook Secret. For full details see [config/config.example.yaml](config/config.example.yaml). Here's a minimal example:

```yaml
repos:
  "your-username/your-repo":
    secret: "your-webhook-secret"
    rules:
      - event: "push"
        branches:
          - main
        actions:
          - type: "shell"
            command: "git pull && ./deploy.sh"
```

### 3. Start the service

```bash
./webhook-listener -config config.yaml -port 9000
```

If you don't need the Web dashboard, you're done. Set your GitHub Webhook URL to `http://your-server:9000/webhook`.

To enable the dashboard, add the `-web` flag:

```bash
./webhook-listener -config config.yaml -port 9000 -web -user admin -pass your-password
```

Then visit `http://your-server:9000/web`.

## CLI Options

| Option | Default | Description |
| --- | --- | --- |
| `-port` | `9000` | HTTP server port |
| `-config` | `config.yaml` | Path to configuration file |
| `-web` | `false` | Enable Web dashboard (at `/web`) |
| `-user` | (empty) | Basic Auth username for Web dashboard |
| `-pass` | (empty) | Basic Auth password for Web dashboard |
| `-workers` | `5` | Maximum concurrent Shell task workers |

## Configuration

For the full configuration format and examples covering three typical use cases, see **[config/config.example.yaml](config/config.example.yaml)**. Below is a field reference.

### Repo Config

| Field | Required | Description |
| --- | --- | --- |
| `name` | No | Display name in the Web dashboard. Defaults to the repository full name |
| `secret` | **Yes** | GitHub Webhook Secret for HMAC-SHA256 signature verification |
| `rules` | **Yes** | List of trigger rules, at least one required |
| `healthcheck` | No | Health check configuration |

### Rule Config

| Field | Required | Description |
| --- | --- | --- |
| `event` | **Yes** | GitHub event type (`push`, `pull_request`, `release`, etc.) |
| `branches` | **Yes** | List of branches to match. Empty list `[]` matches all branches |
| `actions` | **Yes** | Actions to execute on match, at least one required |

### Action Config (type: shell)

| Field | Required | Description |
| --- | --- | --- |
| `type` | **Yes** | Must be `shell` |
| `command` | **Yes** | Shell command to execute. Multi-line text is supported |
| `env` | No | Environment variables, format: `["KEY=VALUE", ...]` |
| `timeout` | No | Timeout in seconds, default `300` |
| `retryCount` | No | Number of retries on failure, default `0` |
| `retryDelay` | No | Seconds between retries, default `0` |
| `workDir` | No | Working directory, defaults to the program's directory |

### Healthcheck Config

| Field | Required | Description |
| --- | --- | --- |
| `url` | **Yes** | URL to probe. The service sends periodic GET requests with a 5-second timeout |
| `interval` | **Yes** | Probe interval in seconds |

> [!NOTE]
> HTTP status codes 200, 301, and 302 are considered healthy. All other status codes or connection failures are treated as unhealthy.

## GitHub Webhook Setup

In your GitHub repository, go to **Settings → Webhooks → Add webhook**:

| Field | Value |
| --- | --- |
| Payload URL | `http://your-server:9000/webhook` |
| Content type | `application/json` |
| Secret | Must match the `secret` field in your config file |
| Events | Select as needed (e.g. `push`, `pull_request`) |

## Web Dashboard

Enabled with the `-web` flag. Access at `/web`. Features:

- **Overview**: Total projects, healthy/unhealthy counts
- **Project Cards**: Latest health status, 24-hour uptime bar chart for each project (hover for details)
- **Deployment Logs**: Last 10 webhook execution records, click to view stdout/stderr output

### Authentication

Basic Auth is recommended:

```bash
./webhook-listener -web -user admin -pass your-password
```

Without credentials configured, the dashboard is accessible without authentication.

### Health Check Endpoint

`GET /healthz` returns `200 OK`. Useful for upstream load balancers or external monitoring tools.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Backend | Go + [Gin](https://github.com/gin-gonic/gin) + [GORM](https://gorm.io/) |
| Database | SQLite (pure-Go driver, zero external dependencies) |
| Frontend | Vue 3 + [Vite](https://vitejs.dev/) + [Tailwind CSS](https://tailwindcss.com/) |
| Config | YAML |

## Directory Structure

```
├── cmd/webhook-listener/main.go   # Entry point
├── config/
│   └── config.example.yaml        # Configuration template
├── internal/
│   ├── bootstrap/app.go           # Dependency injection & wiring
│   ├── config/                    # Config parsing, SQLite initialization
│   ├── handler/                   # HTTP routes & request handling
│   ├── middleware/                # Basic Auth middleware
│   ├── model/                     # GORM models
│   ├── queue/                     # Shell task scheduler, health monitor
│   ├── repository/                # Data access layer
│   ├── service/                   # Business logic layer
│   └── webui/embed.go             # Embedded frontend assets
├── pkg/utils/                     # Utilities (signature, log paths)
├── web/                           # Vue 3 frontend source
└── Makefile
```

## Local Development

```bash
# Backend
go run ./cmd/webhook-listener -config config/config.example.yaml

# Frontend (dev mode with HMR)
cd web && npm install && npm run dev

# Build frontend & cross-compile
make build-all
```

`make build-web` copies `web/dist/` to `internal/webui/dist/`, so the frontend is embedded into the Go binary.

## Notes

- **Network access**: GitHub must be able to reach your `/webhook` endpoint. Make sure your firewall/security group allows the port
- **Signature verification**: The `secret` in your config must match the one in GitHub Webhook settings, or requests will be rejected with HTTP 403
- **HTTPS**: In production, use Nginx/Caddy as a reverse proxy with HTTPS. GitHub displays security warnings for plaintext webhooks
- **Log persistence**: Shell command stdout/stderr are stored under `logs/shell/` in the executable's directory, organized by date
- **SQLite concurrency**: WAL mode is enabled by default, handling concurrent writes well for typical workloads. For extremely high-frequency webhook scenarios, consider migrating to PostgreSQL
