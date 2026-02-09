# Insight Engine – Detailed Documentation

This document explains how the **Insight Engine** works, how to run it end‑to‑end, and what improvements are planned.

## 1. Purpose

The Insight Engine provides predictive insights for CRM contacts (customers). It currently supports:

- **Lead scoring**: Likelihood that a contact converts
- **Churn prediction**: Likelihood that a contact will stop engaging
- **CLV (Customer Lifetime Value)**: Estimated total revenue from a contact

The engine reads features from a **ClickHouse data warehouse** that is populated from the CRM’s transactional Postgres DB via **ETL**.

## 2. Architecture (Runtime Flow)

1. **CRM transactional DB (Postgres)** holds live data: contacts, notes, deals.
2. **ETL** pulls data from Postgres and builds feature tables in ClickHouse.
3. **Insight Engine API (FastAPI)** queries ClickHouse to fetch features.
4. A lightweight model (currently heuristic) calculates scores.
5. The API returns scores to the caller (usually via the Edge API Gateway).

## 3. Folder Structure

```
insight_engine/
  app/
    main.py          # FastAPI service entrypoint
    config.py        # Loads config/default.toml
    warehouse.py     # ClickHouse client helper
    models.py        # Request/response schemas
    lead_scoring.py  # Heuristic scoring logic (lead/churn/clv)

  etl/
    schema.sql       # ClickHouse warehouse tables
    load_contacts_from_postgres.py  # ETL job (incremental + full)
    run_etl.sh       # cron wrapper
    state.json       # (auto-generated) ETL last-run state

  config/default.toml
  docker-compose.yml
  requirements.txt
  README.md
```

## 4. ClickHouse Warehouse

**Tables used:**

- `contacts_features`
- `churn_features`
- `clv_features`

These are rebuilt/updated by ETL. The models read from these tables directly.

### Schema
See: `etl/schema.sql`

## 5. ETL (Extract → Transform → Load)

### How ETL Works

- **Extract:** Reads contacts, notes, and deals from Postgres
- **Transform:** Aggregates feature signals per contact
- **Load:** Writes rows into ClickHouse feature tables

### Incremental ETL

The ETL script supports **incremental mode**:

- Looks for rows updated since last run (`contacts.updated_at`)
- Includes contacts with new notes (`notes.created_at`)
- Includes contacts with updated deals (`deals.updated_at`)
- Deletes existing feature rows for those contacts in ClickHouse
- Inserts updated feature rows

Last run timestamp is stored in `etl/state.json`.

### Full Load

Use `ETL_INCREMENTAL=0` to truncate and reload all feature tables.

## 6. API Endpoints

### 6.1 Lead Score

`POST /v1/lead-score`

Request:
```json
{"lead_id": "CONTACT_ID"}
```

Response:
```json
{"lead_id": "CONTACT_ID", "score": 0.8, "model_version": "v0.1-heuristic"}
```

### 6.2 Churn Score

`POST /v1/churn-score`

Request:
```json
{"contact_id": "CONTACT_ID"}
```

Response:
```json
{"contact_id": "CONTACT_ID", "risk_score": 0.6, "model_version": "v0.1-heuristic"}
```

### 6.3 CLV

`POST /v1/clv`

Request:
```json
{"contact_id": "CONTACT_ID"}
```

Response:
```json
{"contact_id": "CONTACT_ID", "clv": 4200.5, "model_version": "v0.1-heuristic"}
```

## 7. How It Is Used via Edge Gateway

The Edge service proxies these endpoints:

- `/v1/lead-score`
- `/v1/churn-score`
- `/v1/clv`

All requests require `Authorization: Bearer <token>`.

## 8. Configuration

Config file: `config/default.toml`

Sections:
- `[App]`: Host/port for FastAPI
- `[Warehouse]`: ClickHouse connection
- `[Postgres]`: ETL source connection

You can override any field with environment variables in ETL:

```
PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DB
CH_HOST, CH_PORT, CH_USER, CH_PASSWORD, CH_DB
```

## 9. Current Limitations

- Scoring uses heuristic logic, not real ML models
- No feature store or model registry
- No automatic retraining pipeline
- No real-time streaming updates (ETL batch only)
- No evaluation metrics or model monitoring

## 10. Future TODOs

### Data & ETL
1. Add **incremental + watermark** handling for deletes
2. Add **ETL job scheduler** (Airflow/Prefect or systemd timer)
3. Add **feature versioning** and schema evolution

### Modeling
4. Train real models from historical outcomes (lead conversion, churn, CLV)
5. Add **model persistence + versioning** (pickle + metadata or MLflow)
6. Add **offline evaluation metrics** and validation

### API & Serving
7. Add **batch scoring endpoints**
8. Add **caching** for high-frequency requests
9. Add **explainability** outputs (top features)

### Infrastructure
10. Add **monitoring** (Prometheus + Grafana)
11. Add **logging + tracing** for model calls
12. Add **CI/CD pipeline** with automated tests

---

This document reflects the current state of the Insight Engine and the most practical next steps for production readiness.
