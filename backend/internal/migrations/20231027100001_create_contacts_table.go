package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upContactsTable, downContactsTable)
}

func upContactsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS contacts (
			id VARCHAR(36) PRIMARY KEY,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			email VARCHAR(255) UNIQUE,
			phone VARCHAR(50),
			owner_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`)
	return err
}

func downContactsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS contacts;`)
	return err
}
