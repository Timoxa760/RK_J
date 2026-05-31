# Production: VPS (full backend)

Полный бэкенд «Поток» на VPS с Docker, Antigravity (headless), Whisper и TLS через Caddy.

| Компонент | URL / порт |
|-----------|------------|
| API (public) | `https://api.potok.junior.raitokyokai.tech` |
| Gateway (internal) | `127.0.0.1:8000` — только localhost, снаружи через Caddy |
| Antigravity | `http://antigravity-manager:8045` — только Docker-сеть |
| Frontend | Vercel — см. [vercel-deploy.md](./vercel-deploy.md) |

## Требования VPS

- Ubuntu 22.04+ (или аналог)
- Docker Engine + Docker Compose v2
- **8 GB RAM** минимум, **12–16 GB** рекомендуется (Whisper + Kafka + ClickHouse)
- **40+ GB** диск
- Открытые порты: **22**, **80**, **443** (порт **8000 наружу не открывать**)

## DNS

| Запись | Тип | Значение |
|--------|-----|----------|
| `api.potok.junior` | A | IP VPS (например `85.175.100.129`) |

Фронт — CNAME на Vercel (см. vercel-deploy).

## Первичная настройка сервера

```bash
# SSH по ключу (пароль в чат не использовать)
ssh-copy-id user@YOUR_VPS_IP

sudo apt update && sudo apt install -y docker.io docker-compose-plugin caddy
sudo usermod -aG docker $USER
# перелогиниться

sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

## Деплой кода

Из локальной копии ветки **`back`**:

```bash
cd backend
export DEPLOY_HOST=YOUR_VPS_IP
export DEPLOY_USER=your-user
export DEPLOY_PATH=~/potok
./scripts/deploy_server.sh
```

Скрипт: rsync → `docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build` → migrate PG/CH.

Альтернатива: GitHub Actions **Deploy VPS (manual)** на ветке `back` — secrets `VPS_SSH_PRIVATE_KEY`, `VPS_HOST`, `VPS_USER`, `VPS_DEPLOY_PATH`.

Первичная настройка Caddy/ufw на сервере:

```bash
cd ~/potok && bash deploy/bootstrap_vps.sh
```

## `.env` на сервере

```bash
ssh user@YOUR_VPS_IP
cd ~/potok
cp .env.production.example .env
nano .env   # JWT_SECRET, POSTGRES_PASSWORD, CLICKHOUSE_PASSWORD, ANTIGRAVITY_*
```

Обязательные переменные — см. [`.env.production.example`](../../backend/.env.production.example).

После правок:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Caddy (TLS)

```bash
sudo cp ~/potok/deploy/Caddyfile /etc/caddy/Caddyfile
# замените email в Caddyfile
sudo systemctl reload caddy
sudo systemctl status caddy
```

Caddy получит сертификат Let's Encrypt для `api.potok.junior.raitokyokai.tech`.

## Antigravity OAuth (headless)

1. На VPS стек уже поднят с `antigravity-manager` (порт 8045 только внутри Docker).
2. С локальной машины — SSH-туннель к Web UI:

```bash
ssh -L 8045:127.0.0.1:8045 user@YOUR_VPS_IP
# если UI слушает только в контейнере:
ssh user@YOUR_VPS_IP 'docker port backend_antigravity 8045'
# при необходимости: ssh -L 8045:localhost:8045 ...
```

3. Откройте `http://127.0.0.1:8045`, войдите (`ANTIGRAVITY_WEB_PASSWORD`), привяжите Google-аккаунт.
4. Скопируйте API key → `ANTIGRAVITY_API_KEY` в `.env`, перезапустите LLM-сервисы:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d ai-processor credit-service
```

## CORS

В `.env`:

```env
CORS_ALLOWED_ORIGINS=https://potok.junior.raitokyokai.tech
```

Gateway читает `CORS_ALLOWED_ORIGINS` (см. `services/core-api/api-gateway/cors.go`).

## Проверка

```bash
curl -sI https://api.potok.junior.raitokyokai.tech/api/v1/
API_BASE=https://api.potok.junior.raitokyokai.tech ./scripts/smoke_auth_chat.sh
```

## Make (локально / на сервере)

```bash
make up-prod      # compose prod + migrate
make down-prod
```

## Безопасность

- Уникальные `JWT_SECRET`, `ENCRYPTION_KEY`, пароли БД
- Antigravity, Postgres, Redis, Kafka — без public ports в prod compose
- Регулярно обновлять образы; не коммитить `.env`

## Связи

- [vercel-deploy.md](./vercel-deploy.md)
- [environment.md](./environment.md)
- [docker-compose.md](./docker-compose.md)
