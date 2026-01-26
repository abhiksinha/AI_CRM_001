package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upContactTagsTable, downContactTagsTable)
}

func upContactTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS contact_tags (
			contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
			tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (contact_id, tag_id)
		);
	`)
	return err
}

func downContactTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS contact_tags;`)
	return err
}
