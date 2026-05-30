# CI/CD

## Frontend (ветка `front`)

Workflow: `.github/workflows/front.yml`

- Trigger: push/PR в `front`, paths `frontend/**`
- Jobs: `npm ci`, `npm run build`
- Deploy: Vercel (secrets `VERCEL_TOKEN`, `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`)

Remote: `vercel/front` → potok-vercel.

## Backend (ветка `back`)

Workflow: `backend/.github/workflows/back.yml`

- Trigger: push/PR в `back`
- Jobs: `go test ./...`, `docker compose config`
- **Без autodeploy** — только проверки и push в git

Ручной деплой: [backend/docs/backend/DEPLOY.md](../../backend/docs/backend/DEPLOY.md).

## Docs

Изменения API — сначала `docs/api/API_Contract.md`, затем код.
