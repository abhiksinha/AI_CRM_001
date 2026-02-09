CREATE DATABASE IF NOT EXISTS crm_warehouse;

CREATE TABLE IF NOT EXISTS crm_warehouse.contacts_features (
    id String,
    email String,
    phone String,
    notes_count UInt32,
    last_activity_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY id;

CREATE TABLE IF NOT EXISTS crm_warehouse.churn_features (
    contact_id String,
    last_activity_at DateTime,
    notes_count UInt32,
    total_deals UInt32,
    closed_won UInt32,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY contact_id;

CREATE TABLE IF NOT EXISTS crm_warehouse.clv_features (
    contact_id String,
    total_deal_value Float64,
    avg_deal_value Float64,
    total_deals UInt32,
    closed_won UInt32,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY contact_id;
