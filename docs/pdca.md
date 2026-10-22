# PDCA — rabbitmq-cli

## Plan (2026-10-21)

- 安装官方 rabbitmqadmin v2
- 搭 Go 薄封装 + skills + 中文 wiki
- 默认透传；破坏性动词加门禁

## Do

- [x] rabbitmqadmin 2.35.0（本机工具目录，不入库）
- [x] Go skeleton：doctor / skills / admin 透传
- [x] 风险门禁 + 测试（无 `--yes`/`--dry-run` → exit 10）
- [x] skills + wiki
- [x] 私有仓 `github.com/xiaoxiaolin0918/rabbitmq-cli` 已推送

## Check

- `rabbitmqadmin --version` → 2.35.0
- `go test ./...` 绿
- `rabbitmq doctor` 能找到二进制
- 无门禁的 delete → exit 10

## Act / Next

- 可选：把 rabbitmqadmin 目录加入用户 PATH
- 可选：本机有 rabbitmqctl 时多用 `ctl` 透传
- lab 节点只写在本机 `~/.rabbitmqadmin.conf`（永不提交）
- 文档：README + wiki 已补「快速开始 / 门禁 / 运维注意」

## 2026-10-21 真机冒烟（只读 + dry-run）

- **Plan**：对 lab profile（RabbitMQ 4.0.5，rabbitmqadmin 2.35.0）跑 doctor / 只读 / 门禁 / skills，不做真实写入。
- **Do**：先备份 definitions 到 `_backup/<时间戳>/`（gitignore），再跑冒烟脚本。
- **Check**：多数项通过；失败项出在官方工具/服务端（`nodes list`、`plugins list_all` 解析 `/api/nodes` 出错；该版本无 `health_check deprecated_features_in_use`，改用 `deprecated_features list_used`）。
- **Act**（封装修复，已补测）：
  1. 风险门禁：跳过前缀/穿插全局 flag 后再判命令词；补全漏拦写动词（`delete_multiple`、`bind`/`unbind`、`policies set|patch`、`parameters set|clear`、`feature_flags enable`、shovel/federation `declare_*`、`connections close_of_user`、`definitions import_into_vhost`、`config_file` 节点变更、`export`/`get`、`tanzu` 等）。
  2. 门禁 / dry-run 输出中 `--password`/`-p` 脱敏为 `***`。
  3. skills 兜底路径修复（exe 换目录仍能找到 `skills/`）。
  4. `skills read` 拒绝 `..` 与路径分隔符。

## 2026-10-22 真机冒烟（declare → list → delete）

- **Plan**：仅创建带 `rabbitmq-cli-smoke-` 前缀的临时对象，用完删除；操作前 export 备份。
- **Do / Check**：写操作冒烟通过；清理后无前缀残留；definitions 计数与冒烟前一致。
- **Act**：备份与 `_smoke` 产物留在本机，不入库。

## 2026-10-22 文档

- **Plan**：补中文 README 与 wiki；示例 conf 仅 localhost/guest。
- **Do**：重写 README、`docs/wiki/*`，润色 pdca；去掉示例中的真实 lab 地址。
- **Check**：按 README 可 build + `doctor`；门禁与 export/get 原因有专节。
- **Act**：commit 推 `origin/main`（非 force）。
