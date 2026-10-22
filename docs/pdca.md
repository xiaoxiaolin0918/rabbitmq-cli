# PDCA — rabbitmq-cli

## Plan (2026-10-21)

- Install official rabbitmqadmin v2 on iipcloud-xxl
- Scaffold Go thin wrapper + skills + Chinese wiki
- Passthrough first; gates for destructive verbs

## Do

- [x] rabbitmqadmin 2.35.0 → `D:\tools\rabbitmqadmin\rabbitmqadmin.exe`
- [x] Go skeleton build + doctor / skills / admin passthrough
- [x] Risk gates + tests (exit 10 without --yes/--dry-run)
- [x] skills + wiki
- [x] revival AGENT/HANDOFF

## Check

- `rabbitmqadmin --version` → 2.35.0
- `go test ./...` green
- `rabbitmq doctor` finds binary
- blocked delete without gate → exit 10

## Act / Next

- Optional: add `D:\tools\rabbitmqadmin` to user PATH
- Optional: `ctl` passthrough when rabbitmqctl exists
- Wire real lab profile in `~/.rabbitmqadmin.conf` (never commit)
- init git remote when ready (module github.com/xiaoxiaolin0918/rabbitmq-cli)

## 2026-10-21 真机冒烟（只读 + dry-run）

- Plan：针对 192.168.100.41（RabbitMQ 4.0.5，rabbitmqadmin 2.35.0）跑 doctor/只读/门禁/skills，不做真实写入。
- Do：先备份 definitions 到 `_backup/20261021-113010/`（已加入 .gitignore），再跑 `_smoke/smoke.py`。
- Check：68 项通过；3 项失败，都出在官方工具/服务端（`nodes list` 和 `plugins list_all` 解析 /api/nodes 出错 exit 65；2.35 没有 `health_check deprecated_features_in_use` 这个命令，exit 64）。
- Act（封装修复，已补测试）：
  1. risk gate：开头或中间带全局 flag（`--vhost / delete ...`、`queues -V / delete`、`--non-interactive ...`）时能绕过门禁 → 现在先跳过 flag 再判断；补上之前漏掉的写操作动词（delete_multiple、bind/unbind、policies set/patch/declare_*、parameters set/clear、feature_flags enable、shovels/federation declare_*、connections close_of_user、definitions import_into_vhost、config_file add/update/delete_node、*_limits、顶层 export/get、tanzu）。
  2. 门禁/dry-run 输出里的 `--password`/`-p` 值脱敏为 `***`。
  3. skills.go 里的兜底路径含字面 CR（`\r` 被转义成了回车）→ exe 拷到其它目录后报 `skills directory not found`；改用 filepath.Join。
  4. `skills read` 拒绝带路径分隔符或 `..` 的名称。