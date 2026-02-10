# GitHub Webhook Listener

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

A lightweight Go service for receiving GitHub Webhooks and executing Shell commands. It includes an optional built-in web panel for viewing project status, webhook execution logs, and health check results.

**This project has been parsed by Zread. If you need a quick overview of the project, you can click here to view it：[Understand this project](https://zread.ai/zxc7563598/github-webhook-listener)**

| <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00001.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00002.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00003.png"> |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| -                                                                                                    | -                                                                                                    | -                                                                                                    |

## Features

- ​**Webhook Receiver**​: Handles GitHub events such as `push`​, `pull_request`​, `release`, etc.
- ​**Rule Matching**: Configure different rules based on event types and branches to trigger corresponding Shell commands
- ​**Web Panel**: Optional web UI to view runtime overview, webhook logs, and health check status (supports Basic Auth)
- ​**Health Checks**: Optionally configure a URL and interval per repository to periodically probe and record service status

---

## Deployment

### Option 1: Use Prebuilt Releases (Recommended)

Download the binary for your platform from [Releases](https://github.com/zxc7563598/github-webhook-listener/releases), extract it, and run it using the startup commands below.

### Option 2: Build and Run Locally

Requires Go installed locally, and Node.js if you want to build the Web UI.

| Command               | Description                                                 |
| --------------------- | ----------------------------------------------------------- |
| ​`make build`         | Build the executable for the current platform into`bin/`    |
| ​`make run`           | Run using`config.yaml`in the project root (Web UI disabled) |
| ​`make web`           | Same as above, but with the Web UI enabled                  |
| ​`make build-linux`   | Build Linux amd64 (will run`make build-web`first)           |
| ​`make build-darwin`  | Build macOS amd64/arm64                                     |
| ​`make build-windows` | Build Windows amd64                                         |
| ​`make build-all`     | Build for all platforms above                               |
| ​`make build-web`     | Build frontend only and copy to`internal/webui/dist`        |
| ​`make clean`         | Clean build artifacts                                       |

---

## Startup

**Before running, copy** **​`config/config.example.yaml`​**​ **to** **​`config.yaml`​**​ **in the project root and modify it as needed (see configuration details below).**

```bash
./webhook-listener [options]
```

| Option     | Default     | Description                                                       |
| ---------- | ----------- | ----------------------------------------------------------------- |
| ​`-port`   | 9000        | HTTP service listening port                                       |
| ​`-config` | config.yaml | Path to configuration file                                        |
| ​`-web`    | false       | Enable Web UI (accessible at`/web`)                               |
| ​`-user`   | (empty)     | Basic Auth username for Web UI (recommended when`-web`is enabled) |
| ​`-pass`   | (empty)     | Basic Auth password for Web UI                                    |

**Examples:**

```bash
# Webhook only, port 9000, using config.yaml in current directory
./webhook-listener -config config.yaml -port 9000

# Enable Web UI with Basic Auth
./webhook-listener -config config.yaml -port 9000 -web -user admin -pass your-password
```

Set the GitHub Webhook callback URL to: ​`http(s)://your-domain-or-ip:port/webhook`​(e.g. `https://example.com:9000/webhook`)

---

## Configuration

The configuration file uses YAML format. See `config/config.example.yaml` for reference. High-level structure:

### Top Level: `repos`

- key: Full repository name in the format `owner/repo`​ (e.g. `your-username/your-repo`)
- value: Configuration object for that repository

---

### Repository Configuration

| Field          | Required | Description                                    |
| -------------- | -------- | ---------------------------------------------- |
| ​`name`        | No       | Display name shown in the Web UI               |
| ​`secret`      | **Yes**  | GitHub Webhook Secret for signature validation |
| ​`rules`       | **Yes**  | List of rules, at least one required           |
| ​`healthcheck` | No       | Health check configuration (see below)         |

### Rules `rules[]`

| Field       | Required | Description                                           |
| ----------- | -------- | ----------------------------------------------------- |
| ​`event`    | **Yes**  | Event type, such as`push`​,`pull_request`​,`release`  |
| ​`branches` | **Yes**  | List of branches; empty array`[]`matches all branches |
| ​`actions`  | **Yes**  | List of actions, at least one required                |

### Actions `actions[]`​ (currently supports `type: shell`)

| Field         | Required | Description                                                          |
| ------------- | -------- | -------------------------------------------------------------------- |
| ​`type`       | **Yes**  | Must be`shell`                                                       |
| ​`command`    | **Yes**  | Shell command to execute                                             |
| ​`env`        | No       | Environment variables, e.g.`["MY_VAR=hello"]`                        |
| ​`timeout`    | No       | Timeout in seconds, default is 300                                   |
| ​`retryCount` | No       | Number of retries on failure, default is 0                           |
| ​`retryDelay` | No       | Delay between retries (seconds), default is 0                        |
| ​`workDir`    | No       | Working directory for the command; defaults to the program directory |

### Health Check `healthcheck` (Optional)

| Field       | Required | Description                                                                      |
| ----------- | -------- | -------------------------------------------------------------------------------- |
| ​`url`      | **Yes**  | URL to probe via GET with a 5-second timeout; 200/301/302 are considered healthy |
| ​`interval` | **Yes**  | Probe interval in seconds                                                        |

Example configuration snippet:

```yaml
repos:
  # Example repository configuration
  "your-username/your-repo":
    # Display name in the Web UI
    name: "project name"
    # GitHub Webhook Secret (configured in repository settings)
    secret: "your-github-webhook-secret-here"
    rules:
      # Rule 1: Trigger on push events to main or master
      - event: "push"
        branches: ["main", "master"]
        actions:
          - type: "shell" ## Required, must be shell
            command: "git pull && ./deploy.sh" ## Required, shell command to execute
            env: ["MY_VAR=hello", "HTTP_PROXY=http://proxy:8080"] ## Optional, environment variables (similar to one-time exports)
            timeout: 300 ## Optional, execution timeout in seconds (default 300)
            retryCount: 0 ## Optional, retry count on failure (default 0)
            retryDelay: 0 ## Optional, delay between retries in seconds (default 0)
            workDir: "/tmp" ## Optional, working directory (default: binary directory)

      # Rule 2: Trigger on pull_request events for any branch
      - event: "pull_request"
        branches: [] # Empty array matches all branches
        actions:
          - type: "shell"
            command: "echo 'Pull request event received'"

    healthcheck: ## Optional health check
      url: "https://example.com/health" ## A GET request will be sent every interval with a 5s timeout; 200/301/302 are considered healthy
      interval: 30 # Interval in seconds
```

---

## GitHub Webhook Setup Guide

In your GitHub repository, go to:

​`Settings → Webhooks → Add webhook`

Configuration:

- ​**Payload URL**​: `http://your-server:9000/webhook`
- ​**Content type**​: `application/json`
- ​**Secret**: Same value as in the configuration file
- ​**Events**​: Select as needed, e.g. `push`
