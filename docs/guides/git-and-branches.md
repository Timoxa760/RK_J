# Git, ветки и коммиты

## Layout

```
RK_J/                 ← ветка front
  frontend/
  docs/
  .github/workflows/front.yml

RK_J/backend/         ← worktree, ветка back
  services/
  .github/workflows/back.yml
```

```bash
git worktree list
# RK_J         [front]
# RK_J/backend [back]
```

## Коммиты

Формат **Conventional Commits** (Google / Angular):

```
type(scope): subject
```

| type | Пример |
|------|--------|
| feat | `feat(advisor): POST /ai/chat` |
| fix | `fix(credits): empty dashboard without scans` |
| docs | `docs(api): profile skip flags` |
| ci | `ci(front): vercel deploy workflow` |
| refactor | `refactor(profile): file store` |
| test | `test(credits): scan handler` |
| chore | `chore(deps): bump go version` |

**Scope:** `front`, `advisor`, `api`, `ai-processor`, `credits`, `docs`, `ci`.

## Куда коммитить

| Изменения | cwd | Ветка | push |
|-----------|-----|-------|------|
| frontend, docs | корень RK_J | front | `origin front`, Vercel |
| backend | `backend/` | back | `origin back`, CI only |

## Запреты

- Не коммитить `backend/` из корня (front)
- Не коммитить `frontend/` из worktree back
