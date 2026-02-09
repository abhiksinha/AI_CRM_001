import json
import os
from datetime import datetime, timezone
from pathlib import Path
import tomllib

import clickhouse_connect
import psycopg

STATE_FILE = Path(__file__).with_name("state.json")


def to_datetime(ts_seconds: int) -> datetime:
    return datetime.fromtimestamp(ts_seconds, tz=timezone.utc)


def load_state() -> int:
    if not STATE_FILE.exists():
        return 0
    try:
        data = json.loads(STATE_FILE.read_text())
        return int(data.get("last_run", 0))
    except Exception:
        return 0


def save_state(last_run: int) -> None:
    STATE_FILE.write_text(json.dumps({"last_run": int(last_run)}))


def main():
    incremental = os.getenv("ETL_INCREMENTAL", "1") == "1"

    config_path = os.getenv("INSIGHT_CONFIG", "./config/default.toml")
    cfg = {}
    try:
        with open(config_path, "rb") as f:
            cfg = tomllib.load(f)
    except FileNotFoundError:
        cfg = {}

    pg_cfg = cfg.get("Postgres", {})
    ch_cfg = cfg.get("Warehouse", {})

    pg_host = os.getenv("PG_HOST", pg_cfg.get("Host", "127.0.0.1"))
    pg_port = int(os.getenv("PG_PORT", pg_cfg.get("Port", "5432")))
    pg_user = os.getenv("PG_USER", pg_cfg.get("User", "postgres"))
    pg_password = os.getenv("PG_PASSWORD", pg_cfg.get("Password", "postgres"))
    pg_db = os.getenv("PG_DB", pg_cfg.get("Database", "crm"))

    ch_host = os.getenv("CH_HOST", ch_cfg.get("Host", "127.0.0.1"))
    ch_port = int(os.getenv("CH_PORT", ch_cfg.get("Port", "8123")))
    ch_user = os.getenv("CH_USER", ch_cfg.get("User", "default"))
    ch_password = os.getenv("CH_PASSWORD", ch_cfg.get("Password", ""))
    ch_db = os.getenv("CH_DB", ch_cfg.get("Database", "crm_warehouse"))

    pg_dsn = f"host={pg_host} port={pg_port} user={pg_user} password={pg_password} dbname={pg_db}"

    last_run = load_state() if incremental else 0

    with psycopg.connect(pg_dsn) as pg_conn:
        with pg_conn.cursor() as cur:
            if incremental and last_run > 0:
                cur.execute(
                    """
                    WITH changed_contacts AS (
                        SELECT id FROM contacts WHERE updated_at > %s
                        UNION
                        SELECT contact_id AS id FROM notes WHERE created_at > %s
                        UNION
                        SELECT contact_id AS id FROM deals WHERE updated_at > %s
                    )
                    SELECT
                        c.id,
                        c.email,
                        c.phone,
                        COALESCE(n.notes_count, 0) AS notes_count,
                        GREATEST(c.updated_at, COALESCE(n.last_note_at, c.updated_at)) AS last_activity_at,
                        c.updated_at,
                        COALESCE(d.total_deals, 0) AS total_deals,
                        COALESCE(d.closed_won, 0) AS closed_won,
                        COALESCE(d.total_deal_value, 0) AS total_deal_value,
                        COALESCE(d.avg_deal_value, 0) AS avg_deal_value
                    FROM contacts c
                    LEFT JOIN (
                        SELECT contact_id, COUNT(*) AS notes_count, MAX(created_at) AS last_note_at
                        FROM notes
                        GROUP BY contact_id
                    ) n ON n.contact_id = c.id
                    LEFT JOIN (
                        SELECT
                            contact_id,
                            COUNT(*) AS total_deals,
                            SUM(value) AS total_deal_value,
                            AVG(value) AS avg_deal_value,
                            SUM(CASE WHEN stage = 'Closed Won' THEN 1 ELSE 0 END) AS closed_won
                        FROM deals
                        GROUP BY contact_id
                    ) d ON d.contact_id = c.id
                    WHERE c.id IN (SELECT id FROM changed_contacts)
                    """,
                    (last_run, last_run, last_run),
                )
            else:
                cur.execute(
                    """
                    SELECT
                        c.id,
                        c.email,
                        c.phone,
                        COALESCE(n.notes_count, 0) AS notes_count,
                        GREATEST(c.updated_at, COALESCE(n.last_note_at, c.updated_at)) AS last_activity_at,
                        c.updated_at,
                        COALESCE(d.total_deals, 0) AS total_deals,
                        COALESCE(d.closed_won, 0) AS closed_won,
                        COALESCE(d.total_deal_value, 0) AS total_deal_value,
                        COALESCE(d.avg_deal_value, 0) AS avg_deal_value
                    FROM contacts c
                    LEFT JOIN (
                        SELECT contact_id, COUNT(*) AS notes_count, MAX(created_at) AS last_note_at
                        FROM notes
                        GROUP BY contact_id
                    ) n ON n.contact_id = c.id
                    LEFT JOIN (
                        SELECT
                            contact_id,
                            COUNT(*) AS total_deals,
                            SUM(value) AS total_deal_value,
                            AVG(value) AS avg_deal_value,
                            SUM(CASE WHEN stage = 'Closed Won' THEN 1 ELSE 0 END) AS closed_won
                        FROM deals
                        GROUP BY contact_id
                    ) d ON d.contact_id = c.id
                    """
                )

            rows = cur.fetchall()

    if incremental and last_run > 0 and not rows:
        print("No changes since last run")
        return

    client = clickhouse_connect.get_client(
        host=ch_host,
        port=ch_port,
        username=ch_user,
        password=ch_password,
        database=ch_db,
    )

    if not incremental:
        client.command("TRUNCATE TABLE IF EXISTS contacts_features")
        client.command("TRUNCATE TABLE IF EXISTS churn_features")
        client.command("TRUNCATE TABLE IF EXISTS clv_features")

    ids = [row[0] for row in rows]
    if incremental and ids:
        id_list = ",".join([f"'{i}'" for i in ids])
        client.command(f"ALTER TABLE contacts_features DELETE WHERE id IN ({id_list})")
        client.command(f"ALTER TABLE churn_features DELETE WHERE contact_id IN ({id_list})")
        client.command(f"ALTER TABLE clv_features DELETE WHERE contact_id IN ({id_list})")

    contacts_data = []
    churn_data = []
    clv_data = []
    max_ts = last_run
    for (
        contact_id,
        email,
        phone,
        notes_count,
        last_activity_at,
        updated_at,
        total_deals,
        closed_won,
        total_deal_value,
        avg_deal_value,
    ) in rows:
        updated_at = int(updated_at)
        max_ts = max(max_ts, updated_at)
        contacts_data.append(
            [
                contact_id,
                email or "",
                phone or "",
                int(notes_count),
                to_datetime(int(last_activity_at)),
                to_datetime(updated_at),
            ]
        )
        churn_data.append(
            [
                contact_id,
                to_datetime(int(last_activity_at)),
                int(notes_count),
                int(total_deals),
                int(closed_won),
                to_datetime(updated_at),
            ]
        )
        clv_data.append(
            [
                contact_id,
                float(total_deal_value or 0),
                float(avg_deal_value or 0),
                int(total_deals),
                int(closed_won),
                to_datetime(updated_at),
            ]
        )

    if contacts_data:
        client.insert(
            "contacts_features",
            contacts_data,
            column_names=["id", "email", "phone", "notes_count", "last_activity_at", "updated_at"],
        )
    if churn_data:
        client.insert(
            "churn_features",
            churn_data,
            column_names=[
                "contact_id",
                "last_activity_at",
                "notes_count",
                "total_deals",
                "closed_won",
                "updated_at",
            ],
        )
    if clv_data:
        client.insert(
            "clv_features",
            clv_data,
            column_names=[
                "contact_id",
                "total_deal_value",
                "avg_deal_value",
                "total_deals",
                "closed_won",
                "updated_at",
            ],
        )

    if incremental:
        save_state(max_ts)

    print(f"Loaded {len(contacts_data)} contacts into ClickHouse")


if __name__ == "__main__":
    main()
