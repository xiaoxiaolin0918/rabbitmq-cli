# rabbitmq-cli 中文 Wiki

薄封装官方 [`rabbitmqadmin` v2](https://www.rabbitmq.com/docs/management-cli)，不重写官方命令面。

## 快速开始

```powershell
cd D:\work\ops\rabbitmq-cli
.\rabbitmq.ps1 doctor
.\rabbitmq.ps1 skills list
# 透传官方：
.\rabbitmq.ps1 admin show overview
.\rabbitmq.ps1 admin queues list
```

## 配置

优先用官方 `~/.rabbitmqadmin.conf` 的 `[node]` / `-N` 别名。密钥不要进 git。

## 写操作门闩

`delete` / `purge` / `close` / `declare` / `import` 等需要 `--dry-run` 或 `--yes`，否则退出码 10。

## 相关

- 英文 Agent 入口：[`../../AGENTS.md`](../../AGENTS.md)
- 官方文档：https://www.rabbitmq.com/docs/management-cli
