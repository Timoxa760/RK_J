---
name: potok-docs
description: >-
  Документация «Поток» (RK_J, ветка docs) с Context7. Активируется при docs/,
  NAVI.md, Potok_plan.md, OpenAPI, архитектуре, фичах, БД, фазах хакатона.
---

# Поток — документация + Context7

## Карта проекта

Начинай с **[NAVI.md](../../NAVI.md)** — где что искать в репо.

## Когда применять

- Обновление `docs/`, `Potok_plan.md`, `NAVI.md`
- Согласование с ветками `back` / `front` (не ломать порты и структуру)
- OpenAPI, PostgreSQL/ClickHouse, deployment

## Context7

1. Реестр: `.cursor/context7-libraries.json`
2. MCP **context7**: `resolve-library-id` → `query-docs` (до 3 раз)
3. Список id: [libraries.md](./libraries.md)

## Правила

- Язык docs: русский; код — английский
- Раздел `## Связи` в документах компонентов
- `Potok_plan.md` — только оглавление + ссылки
- Social / auction — **гипотеза**, не core demo
- Онбординг — отдельная страница `/onboarding` на `front`
- **API:** правки в `docs/api/API_Contract.md`, `typescript-types.md`, затем `openapi.yaml`

## Связи

- **Зависит от**: MCP context7
- **Связанные документы**: `README.md`, `docs/guides/context7.md`
