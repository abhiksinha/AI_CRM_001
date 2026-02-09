# AI Service (Lead Scoring)

This service provides lead scoring predictions backed by a ClickHouse warehouse.

## Services

- **ClickHouse**: analytics warehouse (`docker-compose.yml`)
- **FastAPI**: prediction API (`/v1/lead-score`, `/v1/churn-score`, `/v1/clv`)

## Quick Start

### 1. Start ClickHouse

```bash
docker compose up -d
```

### 2. Initialize schema

Option A (HTTP):

```bash
curl -u default:clickhouse --data-binary @etl/schema.sql \
  'http://127.0.0.1:8123/?multiquery=1'
```

Option B (client):

```bash
sudo docker exec -i ai_clickhouse clickhouse-client --multiquery < etl/schema.sql
```

### 3. Load data (ETL from Postgres)

Full load (truncate + reload):

```bash
ETL_INCREMENTAL=0 python etl/load_contacts_from_postgres.py
```

Incremental (default):

```bash
python etl/load_contacts_from_postgres.py
```

Configuration is read from `config/default.toml` (can be overridden by env vars).

Override example:

```bash
PG_HOST=127.0.0.1 PG_PORT=5432 PG_USER=postgres PG_PASSWORD=postgres PG_DB=crm \
CH_HOST=127.0.0.1 CH_PORT=8123 CH_USER=default CH_PASSWORD='clickhouse' CH_DB=crm_warehouse \
python etl/load_contacts_from_postgres.py
```

### 4. Run API

```bash
uvicorn app.main:app --reload --port 8090
```

## APIs

### Lead Score

**Endpoint:** `POST /v1/lead-score`

```bash
curl --location 'http://127.0.0.1:8090/v1/lead-score' \
--header 'Content-Type: application/json' \
--data '{
  "lead_id": "LEAD_001"
}'
```

### Churn Score

**Endpoint:** `POST /v1/churn-score`

```bash
curl --location 'http://127.0.0.1:8090/v1/churn-score' \
--header 'Content-Type: application/json' \
--data '{
  "contact_id": "LEAD_001"
}'
```

### CLV

**Endpoint:** `POST /v1/clv`

```bash
curl --location 'http://127.0.0.1:8090/v1/clv' \
--header 'Content-Type: application/json' \
--data '{
  "contact_id": "LEAD_001"
}'
```

## Notes

- `contacts_features`, `churn_features`, and `clv_features` are feature tables.
- The ETL script supports incremental loads (uses `updated_at`, notes, and deals changes).
- State is stored in `etl/state.json`.
