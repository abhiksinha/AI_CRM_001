#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs"
PID_DIR="$ROOT_DIR/.pids"

mkdir -p "$LOG_DIR" "$PID_DIR"

wait_for_port() {
  local host="$1"
  local port="$2"
  local name="$3"
  local retries=30
  local wait=1

  echo "Waiting for $name on ${host}:${port}..."
  for _ in $(seq 1 "$retries"); do
    if (echo >"/dev/tcp/${host}/${port}") >/dev/null 2>&1; then
      echo "$name is up."
      return 0
    fi
    sleep "$wait"
  done

  echo "Timed out waiting for $name on ${host}:${port}."
  return 1
}

if command -v docker >/dev/null 2>&1; then
  echo "Starting dependencies (Postgres, Redis, ClickHouse) via docker compose..."
  DOCKER_COMPOSE="docker compose"
  if ! docker info >/dev/null 2>&1; then
    DOCKER_COMPOSE="sudo docker compose"
  fi
  $DOCKER_COMPOSE -f "$ROOT_DIR/docker-compose.yml" up -d
else
  echo "Docker not found. Ensure Postgres, Redis, and ClickHouse are running."
fi

wait_for_port 127.0.0.1 5432 "Postgres"
wait_for_port 127.0.0.1 6379 "Redis"
wait_for_port 127.0.0.1 8123 "ClickHouse"

run_migrations() {
  local attempts=5
  local delay=2

  for i in $(seq 1 "$attempts"); do
    echo "Running backend migrations (attempt $i/$attempts)..."
    if (cd "$ROOT_DIR/backend" && go run ./cmd/migrations up); then
      return 0
    fi
    echo "Migration attempt $i failed. Waiting ${delay}s before retry..."
    sleep "$delay"
    delay=$((delay + 2))
  done
  return 1
}

if command -v go >/dev/null 2>&1; then
  if ! run_migrations; then
    echo "Failed to run migrations after retries."
    exit 1
  fi
else
  echo "Go not found. Cannot run backend migrations."
  exit 1
fi

init_clickhouse_schema() {
  if command -v curl >/dev/null 2>&1; then
    echo "Initializing ClickHouse schema via HTTP..."
    if curl -sSf -u default:clickhouse --data-binary @"$ROOT_DIR/insight_engine/etl/schema.sql" \
      "http://127.0.0.1:8123/?multiquery=1" >/dev/null; then
      return 0
    fi
    echo "HTTP init failed, falling back to docker exec..."
  fi

  if command -v docker >/dev/null 2>&1; then
    echo "Initializing ClickHouse schema via docker exec..."
    if docker info >/dev/null 2>&1; then
      docker exec -i ai_clickhouse clickhouse-client --multiquery < "$ROOT_DIR/insight_engine/etl/schema.sql"
    else
      sudo docker exec -i ai_clickhouse clickhouse-client --multiquery < "$ROOT_DIR/insight_engine/etl/schema.sql"
    fi
    return 0
  fi

  return 1
}

if ! init_clickhouse_schema; then
  echo "Failed to initialize ClickHouse schema. Ensure curl or docker is available."
  exit 1
fi

PYTHON_BIN=""
if command -v python3 >/dev/null 2>&1; then
  PYTHON_BIN="python3"
elif command -v python >/dev/null 2>&1; then
  PYTHON_BIN="python"
else
  echo "Python not found. Cannot run Insight Engine ETL."
  exit 1
fi

echo "Running Insight Engine ETL (full load)..."
(
  cd "$ROOT_DIR/insight_engine" &&
    ETL_INCREMENTAL=0 "$PYTHON_BIN" etl/load_contacts_from_postgres.py
)

start_process() {
  local name="$1"
  shift
  local logfile="$LOG_DIR/$name.log"
  local pidfile="$PID_DIR/$name.pid"

  if [ -f "$pidfile" ] && kill -0 "$(cat "$pidfile")" >/dev/null 2>&1; then
    echo "$name already running (pid $(cat "$pidfile"))."
    return
  fi

  echo "Starting $name..."
  nohup "$@" > "$logfile" 2>&1 &
  echo $! > "$pidfile"
}

start_process "backend" bash -lc "cd '$ROOT_DIR/backend' && go run ./cmd/api"
start_process "edge" bash -lc "cd '$ROOT_DIR/edge' && go run ./cmd/api"
start_process "insight_engine" bash -lc "cd '$ROOT_DIR/insight_engine' && uvicorn app.main:app --reload --port 8090"

if [ ! -d "$ROOT_DIR/frontend/node_modules" ]; then
  echo "Installing frontend dependencies..."
  (cd "$ROOT_DIR/frontend" && npm install)
fi
start_process "frontend" bash -lc "cd '$ROOT_DIR/frontend' && npm run dev"

cat <<MSG

All services started.
Logs: $LOG_DIR
PIDs: $PID_DIR

Frontend: http://localhost:5173
Edge API: http://127.0.0.1:8081
Backend:  http://127.0.0.1:8080
Insight:  http://127.0.0.1:8090

To stop: run scripts/stop-dev.sh
MSG
