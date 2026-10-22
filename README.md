# rabbitmq-cli

RabbitMQ ops CLI: **thin wrapper** around official [`rabbitmqadmin` v2](https://www.rabbitmq.com/docs/management-cli) (and `rabbitmqctl` when present).

Same habits as nacos-cli / plumelog-cli: profile/node alias, `--yes` / `--dry-run`, Trae-style skills, Chinese wiki.

Module: `github.com/xiaoxiaolin0918/rabbitmq-cli`

## Prepare

1. Official binary (already on this machine): `D:\tools\rabbitmqadmin\rabbitmqadmin.exe`
2. Optional: copy `.env.example` → `.env` (never commit)
3. Prefer official config `~/.rabbitmqadmin.conf` + `-N` node alias

```powershell
# example official config snippet (do not commit real passwords)
# [default]
# host = 127.0.0.1
# port = 15672
# username = guest
# password = guest
```

## Build

```powershell
cd D:\work\ops\rabbitmq-cli
$env:PATH = "D:\tools\go\bin;" + $env:PATH
$env:GOPROXY = "goproxy.cn,direct"
go test ./...
go build -o rabbitmq.exe .
.\rabbitmq.exe doctor
.\rabbitmq.ps1 admin show overview --node default
```

## Gates

High-risk official subcommands (`delete`, `purge`, `close`, `declare`, `import`, …) need `--dry-run` or `--yes`/`-y`. Missing → exit 10.

## Docs

- Agent entry: [`AGENTS.md`](AGENTS.md)
- PDCA: [`docs/pdca.md`](docs/pdca.md)
- Wiki: [`docs/wiki/Home.md`](docs/wiki/Home.md)
