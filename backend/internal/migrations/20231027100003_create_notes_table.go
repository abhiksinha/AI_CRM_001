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
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
			deal_id UUID REFERENCES deals(id) ON DELETE CASCADE,
			author_id UUID REFERENCES users(id) ON DELETE SET NULL,
			content TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func downNotesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS notes;`)
	return err
}
