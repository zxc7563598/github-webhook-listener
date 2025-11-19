# GitHub Webhook Listener

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

一个用于接收 GitHub Webhook 并执行 Shell 命令的轻量级 Go 服务。

**本项目已经经由 Zread 解析完成，如果需要快速了解项目，可以点击此处进行查看：[了解本项目](https://zread.ai/zxc7563598/github-webhook-listener)**

---

## 功能特性

- **安全验证**：支持 GitHub Webhook 签名验证（SHA256）
- **事件过滤**：按事件类型和分支匹配规则
- **Shell 执行**：支持执行 Shell 命令，带超时保护（默认 5 分钟）
- **优雅关闭**：服务退出时等待当前任务执行完成
- **健康检查**：提供 `/health` 端点
- **配置验证**：启动时校验配置文件的有效性
- **请求限制**：限制请求体大小（10MB）并设置合理的超时

---

## 快速开始

### 1. 构建

```bash
go build -o webhook-listener ./cmd/webhook-listener
```

### 2. 创建配置文件

```bash
cp config/config.example.yaml config.yaml
```

编辑 `config.yaml`：

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

### 3. 启动服务

```bash
./webhook-listener -port 9000 -config config.yaml
```

或使用默认配置：

```bash
./webhook-listener
```

默认端口为 `9000`，默认配置文件为 `config.yaml`。

---

## 配置说明

### 配置格式

```yaml
repos:
  "owner/repo":
    secret: "GitHub Webhook Secret" # 必填
    rules:
      - event: "push"
        branches: ["main"]
        actions:
          - type: "shell"
            command: "echo deploy"
```

### 支持事件类型

- ​`push`​
- ​`pull_request`​
- ​`release`​
- 其他 GitHub Webhook 事件

### 分支匹配规则

- ​`branches` 为空或未填写：匹配所有分支
- 指定了分支列表：仅匹配列表内的分支

---

## GitHub Webhook 配置指南

在仓库中进入：

​`Settings → Webhooks → Add webhook`​

配置内容：

- **Payload URL**：`http://your-server:9000/webhook`​
- **Content type**：`application/json`​
- **Secret**：与配置文件一致
- **Events**：根据需求选择，如 `push`​

---

## API 端点

### POST /webhook

接收 GitHub Webhook 请求。

请求头：

- ​`X-GitHub-Event`​
- ​`X-Hub-Signature-256`​

响应状态：

- ​`200 OK`：处理成功
- ​`400 Bad Request`：请求格式错误
- ​`401 Unauthorized`：签名错误
- ​`404 Not Found`：仓库未配置

### GET /health

健康检查接口。

返回示例：

```json
{ "status": "ok" }
```

---

## 安全注意事项

1. **Secret 安全**

   - 不要提交包含真实 secret 的配置文件
   - 建议使用环境变量或密钥管理工具
   - 配置文件权限建议为 `600`​

2. **Shell 执行安全**

   - 不执行来源不明的命令
   - 不拼接用户输入
   - 可根据需要启用命令白名单机制

3. **网络安全**

   - 建议使用 Nginx 或 Caddy 配置 HTTPS
   - 可通过防火墙限制来源 IP

4. **权限控制**

   - 使用最小权限运行服务
   - 控制工作目录读写权限

---

## 部署建议

### 使用 systemd

在 `/etc/systemd/system/webhook-listener.service` 中创建：

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

启用服务：

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

---

## 日志格式

服务会输出结构化日志，包括事件、匹配规则、执行命令及耗时。

示例：

```
[webhook] 仓库: owner/repo, 事件: push, 分支: main
[webhook] 仓库 owner/repo 的规则匹配: event=push, branch=main
[action] executing shell: cd /path && git pull
[shell] 输出: Already up to date.
[webhook] 请求处理完成，耗时: 1.234s
```

---

## 故障排查

### 签名验证失败

- 检查 secret 是否一致
- 确认 GitHub 使用的是 `X-Hub-Signature-256`​

### Shell 命令执行失败

- 检查命令的权限、路径
- 查看日志中的错误输出
- 确认工作目录存在

### 仓库未找到

- 检查 `owner/repo` 名称是否准确
- 注意大小写需与 GitHub 一致

---

## 项目结构

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
