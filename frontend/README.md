# CRM Frontend

React (Vite) frontend for the CRM system. It authenticates via the Edge service, stores the session token in cookies, and exposes core workflows (contacts, deals, tasks, notes, and insight scores).

## Prerequisites

- Node.js 18+ recommended
- Edge service running (default `http://127.0.0.1:8081`)
- Backend and Insight Engine services running behind Edge

## Configure

You can override the Edge base URL with an environment variable:

```bash
export VITE_EDGE_BASE_URL="http://127.0.0.1:8081"
```

## Install

```bash
cd /home/abhishek/GolandProjects/AI_CRM_001/frontend
npm install
```

## Run (Dev)

```bash
npm run dev
```

Open the URL shown by Vite (usually `http://localhost:5173`).

## Run in Docker

From repo root:

```bash
docker compose up --build
```

## Docker Permission Note

If you see Docker permission errors while running the full stack, add your user to the `docker` group and re-login:

```bash
sudo usermod -aG docker $USER
```

## Build

```bash
npm run build
npm run preview
```

## Auth Notes

- Login uses `POST /v1/login`.
- The session token is stored in cookies: `crm_auth_token` and `crm_user_id`.
- On load, the app validates the session by calling `GET /v1/users/{user_id}`.
- Logout calls `POST /v1/logout` and clears cookies.

## Useful Endpoints (via Edge)

- `POST /v1/users` (signup)
- `GET /v1/contacts` / `POST /v1/contacts`
- `GET /v1/deals` / `POST /v1/deals`
- `POST /v1/deals/{dealID}/tasks`
- `POST /v1/contacts/{contactID}/notes`
- `POST /v1/lead-score` / `POST /v1/churn-score` / `POST /v1/clv`
