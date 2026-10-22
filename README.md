# rabbitmq-cli

RabbitMQ 运维 CLI：**官方 [`rabbitmqadmin` v2](https://www.rabbitmq.com/docs/management-cli) 的薄封装**（可选透传 `rabbitmqctl`）。

不重写官方命令面；习惯对齐 nacos-cli / plumelog-cli：profile/`-N` 节点别名、写操作 `--yes` / `--dry-run` 门禁、Trae 风格 skills、中文 wiki。

模块名：`github.com/xiaoxiaolin0918/rabbitmq-cli`

## 它能做什么 / 不能做什么

| 能 | 不能（需另接） |
|----|----------------|
| doctor、overview/list、健康检查等只读快照 | 实时曲线、长期时序落库 |
| 透传官方 `rabbitmqadmin` 全命令面 | 替代 Prometheus / Grafana 监控栈 |
| 写操作门禁（`--yes` / `--dry-run`，缺则 exit 10） | 常驻告警推送（可再包定时 routine） |
| `skills list/read` 给 Agent 用 | 改写官方 CLI 语义 |

监控边界：定位是**运维快照 / 巡检**，不是监控面板。常驻指标请用 RabbitMQ Prometheus 插件 + Grafana。

## 依赖

1. 官方二进制 `rabbitmqadmin` v2（Windows 示例：`D:\tools\rabbitmqadmin\rabbitmqadmin.exe`，或放进 `PATH`）
2. Go 1.24+（仅编译本仓库时需要）
3. 可选：本机已装 RabbitMQ 服务端工具时可用 `rabbitmq ctl ...` 透传

## 安装 / 构建

```powershell
git clone https://github.com/xiaoxiaolin0918/rabbitmq-cli.git
cd rabbitmq-cli

# 按本机 Go 路径调整
$env:PATH = "D:\tools\go\bin;" + $env:PATH
$env:GOPROXY = "goproxy.cn,direct"

go test ./...
go build -o rabbitmq.exe .

# 或用旁路脚本（没有 exe 时会自动 go build）
.\rabbitmq.ps1 doctor
```

首次跑通标准：`.\rabbitmq.exe doctor`（或 `.\rabbitmq.ps1 doctor`）能找到 `rabbitmqadmin` 并打印版本；若尚未配置，会提示复制 `rabbitmqadmin.conf.example`。

## Profile 配置

优先用官方配置文件（不要另造一套 profile 格式）：

1. 复制仓库内示例到用户目录：
   - 源：[`rabbitmqadmin.conf.example`](rabbitmqadmin.conf.example)
   - 目标：`%USERPROFILE%\.rabbitmqadmin.conf`（即 `~/.rabbitmqadmin.conf`）
2. 示例仅含 **localhost / guest**。真实环境地址与密码只写在本机 conf / `.env`，**禁止写进仓库、README、聊天记录**。
3. 用 `-N` / `--node` 选节点别名：

```powershell
.\rabbitmq.ps1 -N default admin show overview
.\rabbitmq.ps1 admin show overview --node default
```

也可用环境变量（见 [`.env.example`](.env.example)），但更推荐官方 conf + `-N`。`.env` 与真实 `rabbitmqadmin.conf` 已在 `.gitignore`。

## 常用命令

```powershell
# 环境自检
.\rabbitmq.ps1 doctor

# 只读总览 / 列表（透传官方）
.\rabbitmq.ps1 admin show overview
.\rabbitmq.ps1 admin show churn
.\rabbitmq.ps1 admin queues list
.\rabbitmq.ps1 admin connections list
.\rabbitmq.ps1 admin consumers list

# 健康检查（无告警时安静通过）
.\rabbitmq.ps1 admin health_check local_alarms
.\rabbitmq.ps1 admin health_check cluster_wide_alarms

# skills（给 Agent / 人读）
.\rabbitmq.ps1 skills list
.\rabbitmq.ps1 skills read rabbitmq-shared

# 可选：本机有 rabbitmqctl 时
.\rabbitmq.ps1 ctl status
```

全局旗标：`--bin`、`--config`/`-c`、`--node`/`-N`、`--dry-run`、`--yes`/`-y`、`--format`（doctor/skills 输出：`json|pretty`）。

## 写操作门禁

高风险透传（含 `declare` / `delete` / `purge` / `close` / `bind` / `unbind` / `import` / `export` / `get` / `publish` 以及 `policies set`、`delete_multiple`、`definitions export|import` 等）必须带：

- `--dry-run`：只打印 would-exec JSON，**不调用**官方二进制；或
- `--yes` / `-y`：确认执行

两者都没有 → **退出码 10**，并打印 `blocked: true`（密码类旗标已脱敏为 `***`）。

```powershell
# 会被拦住（exit 10）
.\rabbitmq.ps1 admin declare queue --name demo

# 先看会执行什么
.\rabbitmq.ps1 admin declare queue --name demo --dry-run

# 确认后写入
.\rabbitmq.ps1 admin declare queue --name demo --yes
```

### 为何 `definitions export` / `export` / `get` 也被门禁？

| 命令 | 原因 |
|------|------|
| `definitions export` / 顶层 `export` | 导出文件含**用户密码哈希**，属于敏感落盘 |
| `get` | 可能 **consume/ack** 消息，不是纯只读 |

冒烟或变更前请先 export 备份到本地时间戳目录，且 **`_backup/` 不入库**（已在 `.gitignore`）。

## 运维注意

- 真机冒烟只造带前缀的临时对象（如 `rabbitmq-cli-smoke-`），用完删除；**不要动已有** vhost / 队列 / 用户 / 权限 / policy。
- 勿提交：`_backup/`、`_smoke/`、`*.exe`、真实 `rabbitmqadmin.conf`、`.env`。
- 官方 2.35.0 上个别子命令（如部分 `nodes list` / `plugins list_all`）可能因服务端响应解析失败；可用 Management HTTP API 或其它 list 命令交叉验证。

## 文档

| 文档 | 内容 |
|------|------|
| [`AGENTS.md`](AGENTS.md) | Agent / 协作者硬规则（英文） |
| [`docs/wiki/Home.md`](docs/wiki/Home.md) | 中文运维 Wiki（快速开始、门禁、运维注意） |
| [`docs/pdca.md`](docs/pdca.md) | PDCA 与冒烟记录 |
| [`skills/`](skills/) | Trae 风格 skills |

## 仓库

Private：https://github.com/xiaoxiaolin0918/rabbitmq-cli
