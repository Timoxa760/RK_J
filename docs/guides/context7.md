# Context7 для проекта «Поток» (RK_J)

## Назначение

Реестр Context7 для актуализации документации.

**Карта проекта:** [NAVI.md](../../NAVI.md)

## Связи

- **MCP**: `context7` в Cursor — [mcp-setup.md](./mcp-setup.md)
- **Skill**: `potok-docs`
- **Аудит**: [../architecture/stack-audit.md](../architecture/stack-audit.md)

## Стек docs ↔ код

| Компонент | Context7 ID | Факт в репо |
|-----------|-------------|-------------|
| Nuxt 4 | `/websites/nuxt_4_x` | front: nuxt 4.3.0 |
| Go | `/golang/go` | back: go 1.25 |
| chi | `/go-chi/docs` | back |
| ClickHouse | `/clickhouse/clickhouse-docs` | 25.12 |
| OpenAPI | `/oai/openapi-specification` | [openapi.yaml](../contracts/openapi.yaml) |

## Как пользоваться

1. MCP context7 — зелёный в Settings
2. В запросе: «potok-docs + context7»
3. До 3× `query-docs` на задачу

## Обновление

При смене версии → `resolve-library-id` → `.cursor/context7-libraries.json`
