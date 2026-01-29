package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upNotesTable, downNotesTable)
}

func upNotesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS notes (
			id VARCHAR(36) PRIMARY KEY,
			contact_id VARCHAR(36) REFERENCES contacts(id) ON DELETE CASCADE,
			deal_id VARCHAR(36) REFERENCES deals(id) ON DELETE CASCADE,
			author_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
			content TEXT,
			created_at BIGINT NOT NULL
		);
	`)
	return err
}

func downNotesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS notes;`)
	return err
}
