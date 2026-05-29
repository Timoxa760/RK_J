# Настройка MCP для проекта «Поток»

> Адаптация [глобального гайда](../../../docs/guides/mcp-setup.md) для RK_J.  
> Skill агента: `potok-docs`, реестр: `.cursor/context7-libraries.json`.

## Зачем

| MCP | Для docs «Поток» |
|-----|-------------------|
| **context7** | Nuxt 4, Go chi, OpenAPI, ClickHouse — без устаревших знаний |
| **github** | RK_J, PR, issues, CI |

## Быстрая проверка

```bash
source ~/.cursor/load-mcp-env.sh 2>/dev/null || true
/home/vr4g/Projects/hack/scripts/check-mcp-env.sh
```

Cursor: **Settings → MCP** — зелёный у `github`, `context7`.

## Конфигурация

| Файл | Назначение |
|------|------------|
| `~/.cursor/mcp.json` | Серверы MCP |
| `~/.cursor/context7-libraries.json` | ID библиотек (глобально) |
| `RK_J/.cursor/context7-libraries.json` | ID для проекта |
| `RK_J/.cursor/skills/potok-docs/` | Skill документации |

## Context7 в чате

```
potok-docs + context7: как описать Bearer JWT в OpenAPI 3.1
```

До 3× `query-docs` на задачу. Реестр: [context7.md](./context7.md).

## GitHub

```bash
gh repo view Timoxa760/RK_J
```

Или MCP github: list branches, PR, checks.

## Связи

- [context7.md](./context7.md)
- [../../NAVI.md](../../NAVI.md)
- Глобальный гайд: `/home/vr4g/Projects/hack/docs/guides/mcp-setup.md`
