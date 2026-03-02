#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$ROOT_DIR/.pids"

if [ -d "$PID_DIR" ]; then
  for pidfile in "$PID_DIR"/*.pid; do
    [ -e "$pidfile" ] || continue
    pid=$(cat "$pidfile")
    if kill -0 "$pid" >/dev/null 2>&1; then
      echo "Stopping $(basename "$pidfile" .pid) (pid $pid)"
      kill "$pid" || true
    fi
    rm -f "$pidfile"
  done
fi

if command -v docker >/dev/null 2>&1; then
  echo "Stopping docker compose dependencies..."
  if docker info >/dev/null 2>&1; then
    docker compose -f "$ROOT_DIR/docker-compose.yml" down
  else
    sudo docker compose -f "$ROOT_DIR/docker-compose.yml" down
  fi
fi
