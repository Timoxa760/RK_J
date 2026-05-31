# Git, ветки и коммиты

## Layout (worktrees)

```
RK_J/                 ← ветка front (frontend + docs копия в front для удобства)
  frontend/

RK_J/backend/         ← worktree, ветка back
  services/
  internal/
  docker-compose.yml

RK_J-docs/            ← worktree, ветка docs (источник правды для документации)
  docs/
  NAVI.md
  CHANGELOG.md
```

```bash
git worktree list
# RK_J         [front]
# RK_J/backend [back]
# RK_J-docs    [docs]
```

## Коммиты

Формат **Conventional Commits**:

```
type(scope): subject
```

| type | Пример |
|------|--------|
| feat | `feat(advisor): POST /ai/chat/stream` |
| fix | `fix(llm): Antigravity OpenAI route` |
| docs | `docs(architecture): advisor system and Antigravity` |
| ci | `ci(front): frontend build workflow` |

## Куда коммитить

| Изменения | cwd | Ветка | push |
|-----------|-----|-------|------|
| `frontend/` | корень RK_J | `front` | `origin front` |
| `backend/` | `backend/` | `back` | `origin back` |
| `docs/`, NAVI, CHANGELOG | `RK_J-docs/` или checkout `docs` | `docs` | `origin docs` |

## Запреты

- Не коммитить `backend/` из корня (front)
- Не коммитить `frontend/` из worktree back
- Не коммитить `.env` с секретами
- Документация — приоритетно в ветке **`docs`**

## Связи

- [NAVI.md](../../NAVI.md)
