# AGENTS.md — rabbitmq-cli

Practical rules for humans and AI agents working on **rabbitmq-cli** (RabbitMQ ops thin wrapper).

Aligns with sibling CLIs **nacos-cli** / **plumelog-cli** / **yunxiao-cli**: English `AGENTS.md` + Chinese `docs/wiki/`.

## Purpose

- Work checkout: `D:\work\ops\rabbitmq-cli` (this tree).
- Bot revival pack (not git): `D:\ai\grok-bot\rabbitmq-cli\` (`AGENT.md` / `HANDOFF.md`).
- **Do not rewrite** official `rabbitmqadmin` / `rabbitmqctl` command surfaces.
- Thin Go + cobra wrapper: profiles/node alias, `--yes` / `--dry-run` gates, doctor, skills, Chinese wiki; default path is **passthrough** to official `rabbitmqadmin`.

## Official tools

| Tool | Role | Default path |
|------|------|--------------|
| `rabbitmqadmin` v2 | HTTP API (queues/exchanges/bindings/users/definitions/…) | `D:\tools\rabbitmqadmin\rabbitmqadmin.exe` |
| `rabbitmqctl` | Node/Erlang (optional; wrap only when present) | on PATH if RabbitMQ server tools installed |

Official config: `~/.rabbitmqadmin.conf` with `-N/--node` aliases. Prefer that over inventing a parallel profile format.

## Methodology (required)

| Practice | Rule |
|----------|------|
| **TDD** | Prefer failing `go test` first for risk gates / binary discovery / passthrough; then implement until green. Run `go test ./...` before claiming done. |
| **DDD** | `internal/domain` (use cases + ports) / `internal/admin` (exec rabbitmqadmin) / `cmd` (cobra). Domain must not import cobra. |
| **PDCA** | Record Plan/Do/Check/Act in `docs/pdca.md`. |

## Directory boundaries

```
rabbitmq-cli/
  AGENTS.md / HANDOFF.md / README.md
  main.go / go.mod / Makefile / rabbitmq.ps1
  cmd/
  internal/admin|risk|output|version
  skills/
  docs/pdca.md / docs/wiki/
```

## Command surface (wrapper)

| Group | Behavior |
|-------|----------|
| `doctor` | Locate rabbitmqadmin, print version, config path, node alias |
| `skills` | `list` `read` `path` |
| `admin` / default passthrough | Forward remaining args to `rabbitmqadmin` |
| `ctl` | Forward to `rabbitmqctl` when available |

Global flags: `--bin`, `--config`/`-c`, `--node`/`-N`, `--dry-run`, `--yes`/`-y`, `--format`.

## High-risk write gates

Destructive passthrough (delete / purge / close / declare overwrite-ish / definitions import) requires `--dry-run` (print would-exec JSON, no call) or `--yes`. Missing gate → exit **10**.

## Secrets — hard rules

Never commit `.env`, `~/.rabbitmqadmin.conf`, passwords, tokens, or `backups/`.

## Build & test

```powershell
cd D:\work\ops\rabbitmq-cli
$env:PATH = "D:\tools\go\bin;" + $env:PATH
$env:GOPROXY = "goproxy.cn,direct"
go test ./...
go build -o rabbitmq.exe .
.\rabbitmq.exe doctor
.\rabbitmq.ps1 --help
```
