# GitHub Webhook Listener

<div align="center">
  <a href="./README.md">English</a>｜<a href="./README.zh-CN.md">简体中文</a>
  <hr width="50%"/>
</div>

用于接收 GitHub Webhook 并执行 Shell 命令的轻量级 Go 服务。内置可选的 Web 面板，支持项目健康监控与 Webhook 执行记录查看。

| <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00001.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00002.png"> | <img src="https://raw.githubusercontent.com/zxc7563598/github-webhook-listener/main/demo/00003.png"> |
| ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |

## 功能特性

- **Webhook 接收**：接收 GitHub 的 push、pull_request、release 等事件，通过 HMAC-SHA256 签名校验
- **规则匹配**：按事件类型、分支灵活配置规则，匹配后自动执行对应 Shell 命令，支持失败重试
- **Web 面板**：可选开启的仪表盘（支持 Basic Auth），查看项目健康状态、24 小时可用率图表、Webhook 执行日志
- **健康监控**：为每个仓库配置定时 HTTP 探测，状态异常一目了然
- **单二进制部署**：Go 编译为单个可执行文件，内置 SQLite，无需外部数据库依赖

## 快速开始

### 1. 获取可执行文件

从 [Releases](https://github.com/zxc7563598/github-webhook-listener/releases) 下载对应平台的二进制，或本地构建：

```bash
make build          # 构建到 bin/
make build-linux    # Linux amd64
make build-darwin   # macOS amd64 + arm64
make build-windows  # Windows amd64
```

> [!NOTE]
> 本地构建需要 Go 1.22+。交叉编译或 Web 面板功能还需 Node.js（用于构建前端）。

### 2. 创建配置文件

```bash
cp config/config.example.yaml config.yaml
```

编辑 `config.yaml`，填入你的仓库和 Webhook Secret。完整说明见 [config/config.example.yaml](config/config.example.yaml)，这里是一个最简示例：

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

### 3. 启动服务

```bash
./webhook-listener -config config.yaml -port 9000
```

如果不需要 Web 面板，到这里就可以用了。GitHub Webhook 地址填 `http://你的服务器:9000/webhook`。

需要仪表盘的，加 `-web` 参数：

```bash
./webhook-listener -config config.yaml -port 9000 -web -user admin -pass your-password
```

然后访问 `http://你的服务器:9000/web`。

## 启动参数

| 参数 | 默认值 | 说明 |
| --- | --- | --- |
| `-port` | `9000` | HTTP 服务监听端口 |
| `-config` | `config.yaml` | 配置文件路径 |
| `-web` | `false` | 开启 Web 面板（访问 `/web`） |
| `-user` | （空） | Web 面板 Basic Auth 用户名 |
| `-pass` | （空） | Web 面板 Basic Auth 密码 |
| `-workers` | `5` | Shell 任务最大并发执行数 |

## 配置文件

配置文件格式和完整示例见 **[config/config.example.yaml](config/config.example.yaml)**，里面包含三种典型场景的配置模板。以下为字段速查。

### Repo 配置

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `name` | 否 | Web 面板中展示的名称，不填则显示仓库全名 |
| `secret` | **是** | GitHub Webhook Secret，用于 HMAC-SHA256 签名校验 |
| `rules` | **是** | 触发规则列表，至少一条 |
| `healthcheck` | 否 | 健康监控配置 |

### Rule 配置

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `event` | **是** | GitHub 事件类型（`push`、`pull_request`、`release` 等） |
| `branches` | **是** | 匹配的分支列表。空列表 `[]` 表示所有分支 |
| `actions` | **是** | 匹配成功后执行的操作列表，至少一项 |

### Action 配置（type: shell）

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `type` | **是** | 固定为 `shell` |
| `command` | **是** | 要执行的 Shell 命令，支持多行文本 |
| `env` | 否 | 环境变量，格式 `["KEY=VALUE", ...]` |
| `timeout` | 否 | 超时时间（秒），默认 `300` |
| `retryCount` | 否 | 失败重试次数，默认 `0` |
| `retryDelay` | 否 | 重试间隔（秒），默认 `0` |
| `workDir` | 否 | 工作目录，默认程序所在目录 |

### Healthcheck 配置

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `url` | **是** | 探测地址。服务定时发送 GET 请求，5 秒超时 |
| `interval` | **是** | 探测间隔（秒） |

> [!NOTE]
> 健康检查将 HTTP 200、301、302 视为正常，其余状态码或连接失败视为异常。

## GitHub Webhook 设置

在 GitHub 仓库中进入 **Settings → Webhooks → Add webhook**：

| 配置项 | 值 |
| --- | --- |
| Payload URL | `http://你的服务器:9000/webhook` |
| Content type | `application/json` |
| Secret | 与配置文件中 `secret` 一致 |
| Events | 按需选择，如 `push`、`pull_request` |

## Web 面板

通过 `-web` 开启，访问 `/web`。面板提供：

- **整体状态**：项目总数、正常/异常数量
- **项目卡片**：每个项目的最新健康状态、24 小时可用率柱状图（悬停显示详情）
- **部署记录**：最近 10 次 Webhook 执行记录，点击可查看 stdout/stderr 日志

### 认证

建议配置 Basic Auth：

```bash
./webhook-listener -web -user admin -pass your-password
```

不配置则面板无需认证即可访问。

### 健康检查端点

服务提供 `GET /healthz`，返回 `200 OK`，可用于上游负载均衡或监控系统的健康探测。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go + [Gin](https://github.com/gin-gonic/gin) + [GORM](https://gorm.io/) |
| 数据库 | SQLite（纯 Go 驱动，零依赖） |
| 前端 | Vue 3 + [Vite](https://vitejs.dev/) + [Tailwind CSS](https://tailwindcss.com/) |
| 配置 | YAML |

## 目录结构

```
├── cmd/webhook-listener/main.go   # 入口
├── config/
│   └── config.example.yaml        # 配置文件模板
├── internal/
│   ├── bootstrap/app.go           # 依赖注入与启动编排
│   ├── config/                    # 配置解析、SQLite 初始化
│   ├── handler/                   # HTTP 路由与请求处理
│   ├── middleware/                # Basic Auth 中间件
│   ├── model/                     # GORM 模型
│   ├── queue/                     # Shell 任务调度器、健康监控器
│   ├── repository/                # 数据访问层
│   ├── service/                   # 业务逻辑层
│   └── webui/embed.go             # 嵌入前端静态文件
├── pkg/utils/                     # 工具函数（签名校验、日志路径）
├── web/                           # Vue 3 前端源码
└── Makefile
```

## 本地开发

```bash
# 后端
go run ./cmd/webhook-listener -config config/config.example.yaml

# 前端（开发模式，支持热更新）
cd web && npm install && npm run dev

# 构建前端并交叉编译
make build-all
```

`make build-web` 会将 `web/dist/` 拷贝到 `internal/webui/dist/`，使前端嵌入 Go 二进制中。

## 注意事项

- **端口对外可达**：GitHub 需要能访问你的 `/webhook` 端点，请确保防火墙/安全组放行对应端口
- **签名校验**：配置文件中的 `secret` 必须与 GitHub Webhook 设置中一致，否则请求会被拒绝（返回 403）
- **HTTPS**：生产环境建议在 Nginx/Caddy 反代后开启 HTTPS，GitHub 对明文传输的 Webhook 有安全提示
- **日志持久化**：Shell 命令的 stdout/stderr 保存在可执行文件所在目录的 `logs/shell/` 下，按日期分目录
- **SQLite 并发**：已启用 WAL 模式，轻量场景下并发写入表现良好。如果是极高频率的 Webhook 场景，可考虑迁移到 PostgreSQL
