package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTasksTable, downTasksTable)
}

func upTasksTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			deal_id UUID REFERENCES deals(id) ON DELETE CASCADE,
			assigned_to_id UUID REFERENCES users(id) ON DELETE SET NULL,
			title VARCHAR(255) NOT NULL,
			due_date TIMESTAMPTZ,
			is_completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

func downTasksTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS tasks;`)
	return err
}
