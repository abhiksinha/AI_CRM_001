# AI CRM

Full-stack CRM with Edge authentication, backend services, and an AI Insight Engine. The frontend is a React (Vite) app that talks to the Edge API.

## Quick Start (All Services in Docker)

```bash
cd /home/abhishek/GolandProjects/AI_CRM_001
docker compose up --build
```

Stop everything:

```bash
docker compose down
```

### Docker Permission Note

If you see Docker permission errors, add your user to the `docker` group and re-login:

```bash
sudo usermod -aG docker $USER
```

## What Gets Started

- Postgres (for backend) via `docker-compose.yml`
- Redis (for Edge auth/session) via `docker-compose.yml`
- ClickHouse (for Insight Engine) via `docker-compose.yml`
- Backend migrations (Goose) on startup
- Insight Engine schema init + ETL full load on startup
- Backend API: `http://127.0.0.1:8080`
- Edge API: `http://127.0.0.1:8081`
- Insight Engine: `http://127.0.0.1:8090`
- Frontend: `http://localhost:5173`

## Manual Run (Per Service, Without Docker)

### Dependencies

```bash
docker compose up -d
```

### Backend (including migrations)

```bash
cd backend
go run ./cmd/migrations up
go run ./cmd/api
```

### Edge

```bash
cd edge
go run ./cmd/api
```

### Insight Engine (schema + ETL + API)

```bash
# schema
curl -u default:clickhouse --data-binary @etl/schema.sql \
  'http://127.0.0.1:8123/?multiquery=1'

# etl
ETL_INCREMENTAL=0 python etl/load_contacts_from_postgres.py

# api
uvicorn app.main:app --reload --port 8090
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Frontend Details

See `frontend/README.md` for UI, auth flow, and endpoints.

## Notes

- The frontend uses cookies `crm_auth_token` and `crm_user_id` for auth.
- Signup uses `POST /v1/users` (Edge proxy to backend).
- Login uses `POST /v1/login`.
