# rabbitmq-cli skills

Trae-style agent skills for the Go thin wrapper.

| Skill | Role |
|-------|------|
| `rabbitmq-shared` | doctor, binary path, config/node, gates |
| `rabbitmq-admin-ops` | passthrough to rabbitmqadmin (list/show/declare/delete/…) |

```powershell
.\rabbitmq.ps1 skills list
.\rabbitmq.ps1 skills read rabbitmq-shared
```
