# GitHub Webhook Listener

一个用于接收 GitHub Webhook 并执行 Shell 命令的 Go 服务。

## 功能特性

- ✅ **安全验证**: 支持 GitHub Webhook 签名验证（SHA256）
- ✅ **事件过滤**: 支持按事件类型和分支过滤
- ✅ **Shell 执行**: 执行配置的 Shell 命令，带超时保护（默认5分钟）
- ✅ **优雅关闭**: 支持优雅关闭，确保正在执行的命令完成
- ✅ **健康检查**: 提供 `/health` 端点用于健康检查
- ✅ **配置验证**: 启动时验证配置文件的有效性
- ✅ **请求限制**: 限制请求体大小（10MB）和超时设置

## 快速开始

### 1. 构建

```bash
go build -o webhook-listener ./cmd/webhook-listener
```

### 2. 配置文件

复制示例配置文件并修改：

```bash
cp config/config.example.yaml config.yaml
```

编辑 `config.yaml`，配置你的仓库和规则：

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

### 3. 运行

```bash
./webhook-listener -port 9000 -config config.yaml
```

或者使用默认配置：

```bash
./webhook-listener
```

默认端口：9000  
默认配置文件：`config.yaml`

## 配置说明

### 配置文件结构

```yaml
repos:
  "仓库全名 (owner/repo)":
    secret: "GitHub Webhook Secret"  # 必填，用于签名验证
    rules:
      - event: "push"                 # GitHub 事件类型
        branches: ["main"]            # 分支列表（空数组表示所有分支）
        actions:
          - type: "shell"             # 操作类型
            command: "echo 'deploy'"  # Shell 命令
```

### 支持的事件类型

- `push`: 代码推送
- `pull_request`: 拉取请求
- `release`: 发布
- 其他 GitHub Webhook 事件类型

### 分支匹配

- 如果 `branches` 为空数组或未指定，则匹配所有分支
- 如果指定了分支列表，只匹配列表中的分支

## GitHub Webhook 配置

在 GitHub 仓库设置中配置 Webhook：

1. 进入仓库 → Settings → Webhooks → Add webhook
2. Payload URL: `http://your-server:9000/webhook`
3. Content type: `application/json`
4. Secret: 与配置文件中的 `secret` 保持一致
5. 选择要触发的事件类型（如 `push`）

## API 端点

### POST /webhook

接收 GitHub Webhook 请求。

**请求头**:
- `X-GitHub-Event`: 事件类型（由 GitHub 自动添加）
- `X-Hub-Signature-256`: 签名（由 GitHub 自动添加）

**响应**:
- `200 OK`: 处理成功
- `400 Bad Request`: 请求格式错误
- `401 Unauthorized`: 签名验证失败
- `404 Not Found`: 仓库未配置

### GET /health

健康检查端点。

**响应**:
```json
{"status":"ok"}
```

## 安全注意事项

1. **Secret 安全**: 
   - 不要在代码仓库中提交包含真实 secret 的配置文件
   - 使用环境变量或密钥管理服务存储 secret
   - 确保配置文件权限设置为 `600`（仅所有者可读）

2. **Shell 命令安全**:
   - 只执行可信的命令
   - 避免执行用户输入的命令
   - 考虑使用命令白名单机制

3. **网络安全**:
   - 建议使用 HTTPS（通过反向代理如 Nginx）
   - 限制访问来源 IP（通过防火墙或反向代理）
   - 使用内网部署，避免暴露到公网

4. **权限控制**:
   - 以最小权限用户运行服务
   - 限制工作目录访问权限

## 部署建议

### 使用 systemd（Linux）

创建 `/etc/systemd/system/webhook-listener.service`:

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

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable webhook-listener
sudo systemctl start webhook-listener
```

### 使用 Nginx 反向代理

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

## 日志

服务会输出结构化日志，包括：
- Webhook 请求信息（仓库、事件、分支）
- 命令执行结果
- 错误信息
- 请求处理时间

日志示例：
```
[webhook] 仓库: owner/repo, 事件: push, 分支: main
[webhook] 仓库 owner/repo 的规则匹配: event=push, branch=main
[action] executing shell: cd /path && git pull
[shell] 输出: Already up to date.
[webhook] 请求处理完成，耗时: 1.234s
```

## 故障排查

1. **签名验证失败**:
   - 检查 GitHub Webhook 配置中的 Secret 是否与配置文件一致
   - 确保使用 `X-Hub-Signature-256` 头（SHA256）

2. **命令执行失败**:
   - 检查命令路径和权限
   - 查看日志中的错误输出
   - 确认工作目录存在且有权限

3. **仓库未找到**:
   - 检查配置文件中的仓库名称格式（必须是 `owner/repo`）
   - 确保仓库名称与 GitHub 中的完全一致

## 开发

### 项目结构

```
.
├── cmd/
│   └── webhook-listener/
│       └── main.go          # 入口文件
├── internal/
│   ├── actions/             # 操作执行
│   │   ├── action.go
│   │   └── shell.go
│   ├── config/              # 配置管理
│   │   └── config.go
│   └── server/              # HTTP 服务器
│       ├── handler.go
│       └── signature.go
├── config/
│   └── config.example.yaml  # 配置示例
└── README.md
```
