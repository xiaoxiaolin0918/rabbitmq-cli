---
name: rabbitmq-admin-ops
version: 0.1.0
description: "Use for RabbitMQ HTTP API ops via rabbitmq-cli admin passthrough: list/show/declare/delete/purge/definitions. Prefer official rabbitmqadmin help. Risk: read | write (needs --yes)."
metadata:
  requires:
    bins: ["rabbitmq", "rabbitmqadmin"]
  cliHelp: "rabbitmqadmin help"
---

# rabbitmq admin ops

Passthrough to official rabbitmqadmin. Explore with `rabbitmqadmin help` / `rabbitmq admin help`.

```powershell
.\rabbitmq.ps1 admin show overview
.\rabbitmq.ps1 admin queues list
.\rabbitmq.ps1 admin delete queue --name demo --dry-run
.\rabbitmq.ps1 admin delete queue --name demo --yes
```
