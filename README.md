# Demonit

 **device monitoring & alerting** stack:

- **Backend** — Go (Gin + GORM + Postgres), hexagonal architecture
- **Frontend** — Nuxt 3 dashboard (Pinia, Tailwind, neo-brutalism UI)
- **Agents** — scripts to heartbeat from Windows / macOS / Linux

Features: device register, heartbeat metrics (CPU/RAM + JSON payload), deadman's switch (auto OFFLINE), API keys, SSE realtime alerts.

## Stack

| Layer | Tech |
|-------|------|
| API | Go 1.22+, Gin, GORM, Zap, validator |
| DB | PostgreSQL |
| UI | Nuxt 3, Vue 3, Pinia, Tailwind |
| Agents | Python / PowerShell / Bash — see [`scripts/`](scripts/README.md) |

## Quick start (Docker)

```bash
cp .env.example .env
docker compose up --build
```

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:3000 |
| API | http://localhost:8080 |
| Health | http://localhost:8080/healthz |
| Postgres | localhost:5432 |

Schema is applied automatically on first DB volume via `backend/migrations/001_init.sql`.  
Agents still run on the host against `http://localhost:8080` (see [`scripts/`](scripts/README.md)).

```bash
docker compose down          # stop
docker compose down -v       # stop + wipe DB volume
```

### Local dev (without Docker)

### 1. Postgres

```sql
CREATE DATABASE demonit;
```

Apply schema (required in production; optional in local if AutoMigrate runs):

```bash
psql -U postgres -d demonit -f backend/migrations/001_init.sql
```

### 2. Backend

```bash
cd backend
cp .env.example .env
# edit DATABASE_DSN, DEVICE_API_KEY, ADMIN_API_KEY
go run ./cmd/server
```

API: `http://localhost:8080` · health: `GET /healthz`

`APP_ENV=development` enables AutoMigrate.  
`APP_ENV=production` skips AutoMigrate — use the SQL migration file.

### 3. Frontend

```bash
cd frontend
cp .env.example .env
# NUXT_ADMIN_API_KEY must match backend ADMIN_API_KEY
npm install
npm run dev
```

Dashboard: `http://localhost:3000`

### 4. Register a device + run an agent

1. Open the dashboard → **Add device** (or `POST /api/v1/devices` with `X-Admin-API-Key`)
2. Copy the device UUID
3. Start an agent:

```bash
# Cross-platform (recommended)
python scripts/device_agent.py --device-id YOUR-UUID

# Windows PowerShell
powershell -ExecutionPolicy Bypass -File scripts/device-agent.ps1 -DeviceId "YOUR-UUID"

# macOS / Linux
chmod +x scripts/device-agent.sh
./scripts/device-agent.sh --device-id YOUR-UUID
```

More detail: [`scripts/README.md`](scripts/README.md)

## Auth model

No user/password login (on purpose for IoT agents).

| Caller | Endpoint | Header |
|--------|----------|--------|
| Device agent | `POST /api/v1/heartbeat` | `X-Device-API-Key` |
| Dashboard register | `POST /api/v1/devices` | `X-Admin-API-Key` (via Nuxt BFF) |
| Dashboard reads | `GET /api/v1/devices*` | public |
| Realtime alerts | `GET /api/v1/events` | public SSE |

## API overview

| Method | Path | Notes |
|--------|------|--------|
| POST | `/api/v1/heartbeat` | device key |
| POST | `/api/v1/devices` | admin key |
| PUT | `/api/v1/devices/:id` | admin key |
| DELETE | `/api/v1/devices/:id` | admin key |
| GET | `/api/v1/devices` | list + latest metrics |
| GET | `/api/v1/devices/:id` | detail |
| GET | `/api/v1/devices/:id/metrics` | history |
| GET | `/api/v1/events` | SSE `device_offline` |
| GET | `/healthz` | health |

## Tests

```bash
cd backend
go test ./...
```

## Project layout

```
backend/     Go API (hexagonal: domain / port / adapter / application)
frontend/    Nuxt 3 dashboard
scripts/     Device heartbeat agents
```

## License

MIT — see [LICENSE](LICENSE).
