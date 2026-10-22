---
name: rabbitmq-shared
version: 0.1.0
description: "Use for rabbitmq-cli setup: doctor, locate rabbitmqadmin, ~/.rabbitmqadmin.conf + --node, --dry-run/--yes gates, skills list/read. Risk: read (doctor) | write (gated passthrough)."
metadata:
  requires:
    bins: ["rabbitmq", "rabbitmqadmin"]
  cliHelp: "rabbitmq --help"
---

# rabbitmq-cli shared

Read this before other `rabbitmq-*` skills. Mainline is **Go thin wrapper** over official rabbitmqadmin v2.

## Do NOT

- Paste real passwords into chat / docs / git
- Commit `.env` or `~/.rabbitmqadmin.conf`
- Run destructive admin verbs without `--dry-run` then approved `--yes`

## Steps

1. `rabbitmq doctor` — confirm binary + version
2. Configure `~/.rabbitmqadmin.conf` node alias (or env `RABBITMQADMIN_*`)
3. Read ops: `rabbitmq admin show overview` / `queues list`
4. Writes: `--dry-run` first, then `--yes`
