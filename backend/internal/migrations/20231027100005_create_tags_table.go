package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTagsTable, downTagsTable)
}

func upTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tags (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL
		);
	`)
	return err
}

func downTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS tags;`)
	return err
}
