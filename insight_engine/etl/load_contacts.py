import os
import clickhouse_connect


def main():
    host = os.getenv("CH_HOST", "127.0.0.1")
    port = int(os.getenv("CH_PORT", "8123"))
    user = os.getenv("CH_USER", "default")
    password = os.getenv("CH_PASSWORD", "")
    database = os.getenv("CH_DB", "crm_warehouse")

    client = clickhouse_connect.get_client(
        host=host,
        port=port,
        username=user,
        password=password,
        database=database,
    )

    # Example ETL stub: insert a demo row (replace with real extraction)
    client.command(
        "INSERT INTO contacts_features (id, email, phone, notes_count, last_activity_at) VALUES",
        [
            ["LEAD_001", "lead@example.com", "+10000000000", 2, "2026-02-01 10:00:00"],
        ],
    )

    print("Inserted demo row into contacts_features")


if __name__ == "__main__":
    main()
