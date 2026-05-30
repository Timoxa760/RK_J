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
- **Без autodeploy и без docker compose** — только проверки и push в git

Ручной деплой: [backend/docs/backend/DEPLOY.md](../../backend/docs/backend/DEPLOY.md).

## Docs

Изменения API — сначала `docs/api/API_Contract.md`, затем код.
