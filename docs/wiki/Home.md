# rabbitmq-cli 中文 Wiki

薄封装官方 [`rabbitmqadmin` v2](https://www.rabbitmq.com/docs/management-cli)，**不重写**官方命令面。默认路径是透传；本仓库只加 doctor、skills、`-N` 节点别名，以及写操作门禁。

入口导航：[Wiki 目录](README.md) · [仓库 README](../../README.md) · [AGENTS.md](../../AGENTS.md) · [PDCA](../pdca.md)

## 快速开始

### 1. 构建

```powershell
cd <checkout>\rabbitmq-cli
$env:PATH = "D:\tools\go\bin;" + $env:PATH   # 按本机调整
$env:GOPROXY = "goproxy.cn,direct"
go test ./...
go build -o rabbitmq.exe .
# 或
.\rabbitmq.ps1 doctor
```

### 2. 配置节点（localhost 示例）

```powershell
Copy-Item .\rabbitmqadmin.conf.example $env:USERPROFILE\.rabbitmqadmin.conf
# 编辑 conf：示例默认 127.0.0.1:15672 / guest。真实 lab 只改本机文件，勿回写仓库。
```

### 3. 验证

```powershell
.\rabbitmq.ps1 doctor
.\rabbitmq.ps1 -N default admin show overview
.\rabbitmq.ps1 admin queues list
.\rabbitmq.ps1 skills list
```

`doctor` 成功标准：找到 `rabbitmqadmin`、打印版本、报告默认 conf 是否存在。

## 门禁与安全

封装在调用官方二进制前做风险判定（见 `internal/risk/gate.go`）：

1. 跳过前缀/穿插的全局旗标（`-N`、`--vhost`、`-u`/`-p` 等），再取命令词 —— 避免 `--vhost / delete ...` 这类写法绕过门禁。
2. 判定为高风险且既无 `--dry-run` 也无 `--yes` → **exit 10**，stdout 为 JSON（`blocked: true`）。
3. `--dry-run`：打印 would-exec（`--password` / `-p` 值显示为 `***`），不执行。
4. `--yes` / `-y`：放行并透传。

### 覆盖的写操作（摘要）

- 顶层：`declare` / `delete` / `purge` / `close` / `import` / `export` / `publish` / `rebalance` / `get`
- 子命令动词及前缀：`bind` / `unbind` / `set` / `patch` / `clear` / `enable*` / `declare_*` / `delete_*` / `close_*` / `import_*` / `add_*` 等
- `definitions export|import`（含 `import_into_vhost`）
- 其它曾漏拦、已补测：`policies set|patch`、`parameters set|clear`、`delete_multiple`、`connections close_of_user`、`feature_flags enable`、shovel/federation 的 `declare_*`、`config_file` 节点增删改等

### 敏感落盘与消费

| 行为 | 门禁原因 | 建议 |
|------|----------|------|
| `definitions export` / `export` | 文件含密码哈希 | 先备份到 `_backup/<时间戳>/`，该目录已 gitignore |
| `get` | 可能消费消息 | 只读排查优先用 Management UI / `queues list` 指标 |
| 聊天/文档/PR | — | 禁止粘贴真实密码、lab 地址 |

`skills read` 只允许读 `skills/` 下技能名，拒绝 `..` 与路径分隔符。

## 运维注意

### 真机冒烟约定

1. **先** `definitions export` 到本地时间戳目录（需 `--yes`）。
2. 顺序：只读 → `--dry-run` → 写操作只造**新**对象。
3. 临时对象统一前缀，例如 `rabbitmq-cli-smoke-`（vhost / queue / exchange / user 皆然）。
4. 用完 `delete`（带 `--yes`）清理；核对无前缀残留、definitions 对象计数与冒烟前一致。
5. **禁止**改已有 vhost、队列、用户、权限、policy。

### 不要入库的路径

已在 `.gitignore`：`_backup/`、`_smoke/`、`*.exe`、`rabbitmqadmin.conf`、`.env`、`backups/`。

### 监控边界

本 CLI 适合：

- 变更前后快照（overview / queues / alarms）
- 人工或定时巡检（exit code + 日志）

不适合单独承担：实时曲线、长期 retention、告警路由 —— 请接 Prometheus 插件或外部 routine。

### 已知官方边角（2.35.0 实机）

- 部分环境下 `nodes list` / `plugins list_all` 解析 `/api/nodes` 可能失败；节点本身可用 HTTP API 或其它 list 验证。
- `health_check deprecated_features_in_use` 在该版本可能不存在；改用 `deprecated_features list_used`。

## 命令面速查

| 封装命令 | 行为 |
|----------|------|
| `doctor` | 定位 rabbitmqadmin、版本、默认 conf、可选 rabbitmqctl |
| `admin …` / 透传 | 剩余参数交给 `rabbitmqadmin` |
| `ctl …` | 本机有 `rabbitmqctl` 时透传 |
| `skills list\|read\|path` | 读仓库内 skills |

官方完整帮助：`rabbitmqadmin help` 或 `.\rabbitmq.ps1 admin help`。

## 相关链接

- 官方 Management CLI：https://www.rabbitmq.com/docs/management-cli
- Skills 说明：[`../../skills/README.md`](../../skills/README.md)
