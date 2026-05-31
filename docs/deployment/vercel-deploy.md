# Production: Vercel (frontend)

Nuxt 4 фронт «Поток» на Vercel с доменом `potok.junior.raitokyokai.tech` и API на отдельном HTTPS-домене.

## Архитектура

```
Browser → https://potok.junior.raitokyokai.tech     (Vercel)
Browser → https://api.potok.junior.raitokyokai.tech (VPS + Caddy)
```

Фронт **не** проксирует API через Vercel rewrites — SSE-чат советника идёт напрямую на API (CORS + TLS). См. `composables/useApi.ts`, `utils/advisorStream.ts`.

## Репозиторий

| Параметр | Значение |
|----------|----------|
| Ветка | `front` |
| Root Directory | `frontend` |
| Framework | Nuxt |
| Build | `npm run build` |
| Install | `npm ci` |

Файлы: [`frontend/vercel.json`](../../frontend/vercel.json), [`frontend/nuxt.config.ts`](../../frontend/nuxt.config.ts) (`nitro.preset: 'vercel'` при `VERCEL=1`).

## Environment Variables (Production)

| Variable | Value |
|----------|-------|
| `NUXT_PUBLIC_API_BASE` | `https://api.potok.junior.raitokyokai.tech` |
| `NUXT_PUBLIC_DEMO_MODE` | `false` |

Шаблон: [`.env.production.example`](../../frontend/.env.production.example).

## Custom domain

1. Vercel → Project → **Domains** → Add `potok.junior.raitokyokai.tech`
2. В DNS зоны `raitokyokai.tech`:

| Запись | Тип | Значение |
|--------|-----|----------|
| `potok.junior` | CNAME | значение из Vercel (например `cname.vercel-dns.com`) |

3. Дождаться выпуска сертификата Vercel (Automatic HTTPS).

## Deploy

- Подключить GitHub/GitLab репозиторий, ветка `front`
- Production deploy — push в `front` или manual deploy из dashboard
- Preview deployments используют те же env или override `NUXT_PUBLIC_API_BASE` на staging API

## E2E checklist

1. `https://potok.junior.raitokyokai.tech` — лендинг
2. `/login` — регистрация/вход (запросы на `api.potok...`)
3. `/advisor` — streaming SSE, без CORS/mixed-content ошибок
4. Голосовой ввод — whisper на VPS через API

Backend smoke:

```bash
API_BASE=https://api.potok.junior.raitokyokai.tech ./backend/scripts/smoke_auth_chat.sh
API_BASE=https://api.potok.junior.raitokyokai.tech ./backend/scripts/prod_e2e_check.sh
```

## Связи

- [vps-deploy.md](./vps-deploy.md)
- [front-quickstart.md](./front-quickstart.md)
- [environment.md](./environment.md)
