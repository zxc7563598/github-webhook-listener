# GitHub Webhook Listener

用于接收 GitHub Webhook 并执行 Shell 命令的轻量级 Go 服务，内置 Web 面板可查看项目运行状态与 Webhook 执行记录。

## 功能特性

- **Webhook 接收**：接收 GitHub 的 push、pull_request、release 等事件
- **规则匹配**：按事件类型、分支配置不同规则，触发对应 Shell 命令
- **Web 面板**：可选开启，查看运行概览、Webhook 日志、健康监控状态（支持 Basic Auth）
- **健康监控**：可选为每个仓库配置 URL 与间隔，定时探测并记录状态

---

## 部署方式

### 方式一：使用 Releases 预构建（推荐）

在 [Releases](https://github.com/zxc7563598/github-webhook-listener/releases) 下载对应平台的二进制，解压后按下方「启动命令」运行即可。

### 方式二：本地构建与运行

需要本地安装 Go 与（若需 Web 面板）Node.js。

| 命令                 | 说明                                            |
| -------------------- | ----------------------------------------------- |
| `make build`         | 构建当前平台可执行文件到 `bin/`                 |
| `make run`           | 使用项目根目录 `config.yaml` 运行（不开启 Web） |
| `make web`           | 同上，并开启 Web 面板                           |
| `make build-linux`   | 构建 Linux amd64（会先执行 `make build-web`）   |
| `make build-darwin`  | 构建 macOS amd64/arm64                          |
| `make build-windows` | 构建 Windows amd64                              |
| `make build-all`     | 构建上述所有平台                                |
| `make build-web`     | 仅构建前端并拷贝到 `internal/webui/dist`        |
| `make clean`         | 清理构建产物                                    |

---

## 启动命令

**运行前需将 `config/config.example.yaml` 复制为项目根目录下的 `config.yaml` 并按要求修改（见下方配置说明）。**

```bash
./webhook-listener [选项]
```

| 参数      | 默认值      | 说明                                                 |
| --------- | ----------- | ---------------------------------------------------- |
| `-port`   | 9000        | HTTP 服务监听端口                                    |
| `-config` | config.yaml | 配置文件路径                                         |
| `-web`    | false       | 是否开启 Web 面板（访问 `/web`）                     |
| `-user`   | （空）      | Web 面板 Basic Auth 用户名（开启 `-web` 时建议设置） |
| `-pass`   | （空）      | Web 面板 Basic Auth 密码                             |

**示例：**

```bash
# 仅 Webhook，端口 9000，使用当前目录 config.yaml
./webhook-listener -config config.yaml -port 9000

# 开启 Web 面板，并设置 Basic Auth
./webhook-listener -config config.yaml -port 9000 -web -user admin -pass your-password
```

GitHub Webhook 回调地址填写：`http(s)://你的域名或IP:端口/webhook`（例如 `https://example.com:9000/webhook`）。

---

## 配置文件说明

配置文件为 YAML，参考 `config/config.example.yaml`。结构概览如下。

### 顶层：`repos`

- key：仓库全名，格式 `owner/repo`（如 `your-username/your-repo`）
- value：该仓库的配置对象

---

### 每个仓库的配置

| 字段          | 必填   | 说明                                |
| ------------- | ------ | ----------------------------------- |
| `name`        | 否     | 在 Web 面板中显示的名称             |
| `secret`      | **是** | GitHub Webhook Secret，用于签名校验 |
| `rules`       | **是** | 规则列表，至少一条                  |
| `healthcheck` | 否     | 健康监控配置（见下）                |

### 规则 `rules[]`

| 字段       | 必填   | 说明                                           |
| ---------- | ------ | ---------------------------------------------- |
| `event`    | **是** | 事件类型，如 `push`、`pull_request`、`release` |
| `branches` | **是** | 分支列表；空数组 `[]` 表示匹配所有分支         |
| `actions`  | **是** | 操作列表，至少一项                             |

### 操作 `actions[]`（当前支持 type: shell）

| 字段         | 必填   | 说明                                |
| ------------ | ------ | ----------------------------------- |
| `type`       | **是** | 固定为 `shell`                      |
| `command`    | **是** | 要执行的 Shell 命令                 |
| `env`        | 否     | 环境变量列表，如 `["MY_VAR=hello"]` |
| `timeout`    | 否     | 超时时间（秒），默认 300            |
| `retryCount` | 否     | 失败后重试次数，默认 0              |
| `retryDelay` | 否     | 重试间隔（秒），默认 0              |
| `workDir`    | 否     | 命令工作目录，默认为程序所在目录    |

### 健康监控 `healthcheck`（可选）

| 字段       | 必填   | 说明                                                |
| ---------- | ------ | --------------------------------------------------- |
| `url`      | **是** | 要探测的地址（GET），5 秒超时；200/301/302 视为成功 |
| `interval` | **是** | 探测间隔（秒）                                      |

配置示例片段：

```yaml
repos:
  # 示例：配置一个仓库
  "your-username/your-repo":
    # 在 Web 控制台展示的名称
    name: "project name"
    # GitHub Webhook Secret（在 GitHub 仓库设置中配置）
    secret: "your-github-webhook-secret-here"
    rules:
      # 规则1: 当 main 或 master 分支有 push 事件时触发
      - event: "push"
        branches: ["main", "master"]
        actions:
          - type: "shell" ## 必须，固定为 shell
            command: "git pull && ./deploy.sh" ## 必须，要执行的 cmd 命令
            env: ["MY_VAR=hello", "HTTP_PROXY=http://proxy:8080"] ## 非必须，环境变量，约等于一次性的 export
            timeout: 300 ## 非必须，允许执行时间，默认 300，单位秒
            retryCount: 0 ## 非必须，失败后重试次数，默认 0
            retryDelay: 0 ## 非必须，失败后间隔多久进行重试，默认 0，单位秒
            workDir: "/tmp" ## 非必须，命令工作目录，默认项目二进制文件目录

      # 规则2: 当任何分支有 pull_request 事件时触发
      - event: "pull_request"
        branches: [] # 空数组表示所有分支
        actions:
          - type: "shell"
            command: "echo 'Pull request event received'"

    healthcheck: ## 非必须，健康监控，配置后可以进行健康监控
      url: "https://example.com/health" ## 健康监控的地址，配置后会按照 interval 的间隔发起限时 5 秒的 GET 请求，200/301/302 视为成功
      interval: 30 # 监控间隔，单位秒
```

---

## GitHub Webhook 配置指南

在仓库中进入：

​`Settings → Webhooks → Add webhook`​

配置内容：

- **Payload URL**：`http://your-server:9000/webhook`​
- **Content type**：`application/json`​
- **Secret**：与配置文件一致
- **Events**：根据需求选择，如 `push`​
