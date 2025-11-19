# GitHub Webhook Listener

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

A lightweight Go service for receiving GitHub Webhooks and executing shell commands.

**This project has been parsed by Zread. If you need a quick overview of the project, you can click here to view it：[Understand this project](https://zread.ai/zxc7563598/github-webhook-listener)**

---

## Features

- **Secure Verification**: Supports GitHub Webhook signature verification (SHA256)
- **Event Filtering**: Filter by event type and branch matching rules
- **Shell Execution**: Supports executing shell commands with timeout protection (default 5 minutes)
- **Graceful Shutdown**: Waits for ongoing tasks to finish when the service exits
- **Health Check**: Provides the `/health` endpoint
- **Configuration Validation**: Validates configuration file on startup
- **Request Limiting**: Limits request body size (10MB) and sets reasonable timeouts

---

## Quick Start

### 1. Build

```bash
go build -o webhook-listener ./cmd/webhook-listener
```

### 2. Create a configuration file

```bash
cp config/config.example.yaml config.yaml
```

Edit `config.yaml`:

```yaml
repos:
  "your-username/your-repo":
    secret: "your-github-webhook-secret"
    rules:
      - event: "push"
        branches: ["main", "master"]
        actions:
          - type: "shell"
            command: "cd /path/to/your/project && git pull && ./deploy.sh"
```

### 3. Start the service

```bash
./webhook-listener -port 9000 -config config.yaml
```

Or use default settings:

```bash
./webhook-listener
```

The default port is `9000`, and the default configuration file is `config.yaml`.

---

## Configuration Guide

### Configuration Format

```yaml
repos:
  "owner/repo":
    secret: "GitHub Webhook Secret" # Required
    rules:
      - event: "push"
        branches: ["main"]
        actions:
          - type: "shell"
            command: "echo deploy"
```

### Supported Event Types

- ​`push`​
- ​`pull_request`​
- ​`release`​
- Other GitHub Webhook events

### Branch Matching Rules

- If `branches` is empty or omitted: match all branches
- If branch list is specified: only match branches in the list

---

## GitHub Webhook Setup Guide

In your repository:

​`Settings → Webhooks → Add webhook`​

Configure:

- **Payload URL**: `http://your-server:9000/webhook`​
- **Content type**: `application/json`​
- **Secret**: Same as in the configuration file
- **Events**: Choose as needed, e.g., `push`​

---

## API Endpoints

### POST /webhook

Receives GitHub Webhook requests.

Headers:

- ​`X-GitHub-Event`​
- ​`X-Hub-Signature-256`​

Response Status:

- ​`200 OK`: Processed successfully
- ​`400 Bad Request`: Invalid request format
- ​`401 Unauthorized`: Signature verification failed
- ​`404 Not Found`: Repository not configured

### GET /health

Health check endpoint.

Example response:

```json
{ "status": "ok" }
```

---

## Security Notes

1. **Secret Security**

   - Do not commit configuration files containing real secrets
   - Recommended to use environment variables or secret managers
   - Set configuration file permissions to `600`​

2. **Shell Execution Safety**

   - Do not run commands from untrusted sources
   - Do not concatenate user input
   - Enable a command whitelist mechanism if needed

3. **Network Security**

   - Recommended to use Nginx or Caddy to enable HTTPS
   - Firewall can be used to restrict source IPs

4. **Permission Control**

   - Run the service with minimal required permissions
   - Limit read/write permissions of the working directory

---

## Deployment Suggestions

### Using systemd

Create `/etc/systemd/system/webhook-listener.service`:

```ini
[Unit]
Description=GitHub Webhook Listener
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/webhook-listener
ExecStart=/path/to/webhook-listener -port 9000 -config /path/to/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable webhook-listener
sudo systemctl start webhook-listener
```

### Using Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /webhook {
        proxy_pass http://localhost:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://localhost:9000;
    }
}
```

---

## Log Format

The service outputs structured logs including events, matching rules, executed commands, and execution duration.

Example:

```
[webhook] Repository: owner/repo, Event: push, Branch: main
[webhook] Rule matched for owner/repo: event=push, branch=main
[action] executing shell: cd /path && git pull
[shell] Output: Already up to date.
[webhook] Request completed, duration: 1.234s
```

---

## Troubleshooting

### Signature verification failed

- Check whether the secret matches
- Confirm GitHub is sending `X-Hub-Signature-256`​

### Shell command execution failed

- Check command permissions and paths
- Check log output for error messages
- Ensure the working directory exists

### Repository not found

- Check whether `owner/repo` is correct
- Note GitHub repository names are case-sensitive

---

## Project Structure

```
.
├── cmd/
│   └── webhook-listener/
│       └── main.go          # Entry file
├── internal/
│   ├── actions/             # Action execution
│   │   ├── action.go
│   │   └── shell.go
│   ├── config/              # Configuration management
│   │   └── config.go
│   └── server/              # HTTP server
│       ├── handler.go
│       └── signature.go
├── config/
│   └── config.example.yaml  # Configuration example
└── README.md
```
