#!/usr/bin/env bash
set -euo pipefail

cd /home/abhishek/GolandProjects/AI_CRM_001/insight_engine
python etl/load_contacts_from_postgres.py
