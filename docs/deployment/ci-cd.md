# CI/CD

## Frontend (ветка `front`)

Workflow: `.github/workflows/front.yml`

- Trigger: push/PR в `front`, paths `frontend/**`
- Jobs: `npm ci`, `npm run build`
- **Без autodeploy** — только проверки и push в git (`origin/front`)

## Backend (ветка `back`)

Workflow: `backend/.github/workflows/back.yml`

- Trigger: push/PR в `back`
- Jobs: `go test ./...`
- **Без autodeploy** — только проверки и push в git

Ручной деплой на VPS: [vps-deploy.md](./vps-deploy.md) (`scripts/deploy_server.sh`, `deploy/bootstrap_vps.sh`).

Опционально: `backend/.github/workflows/deploy-vps.yml` (manual) — secrets `VPS_SSH_PRIVATE_KEY`, `VPS_HOST`, `VPS_USER`.

## Production front (Vercel)

См. [vercel-deploy.md](./vercel-deploy.md) — домен `potok.junior.raitokyokai.tech`, env `NUXT_PUBLIC_API_BASE`.

## Docs

Изменения API — сначала `docs/api/API_Contract.md`, затем код.
