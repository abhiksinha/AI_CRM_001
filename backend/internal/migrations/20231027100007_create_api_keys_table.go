package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upApiKeysTable, downApiKeysTable)
}

func upApiKeysTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS api_keys (
			id VARCHAR(18) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at BIGINT NOT NULL,
			deleted_at BIGINT
		);
	`)
	if err != nil {
		return err
	}

	// This partial unique index ensures a user can only have one active (not deleted) API key.
	_, err = tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX user_active_api_key_idx ON api_keys (user_id) WHERE deleted_at IS NULL;
	`)
	return err
}

func downApiKeysTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS api_keys;`)
	return err
}
