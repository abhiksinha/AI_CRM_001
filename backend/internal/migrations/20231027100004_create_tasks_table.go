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
			id VARCHAR(36) PRIMARY KEY,
			deal_id VARCHAR(36) REFERENCES deals(id) ON DELETE CASCADE,
			assigned_to_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
			title VARCHAR(255) NOT NULL,
			due_date BIGINT,
			is_completed BOOLEAN DEFAULT FALSE,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`)
	return err
}

func downTasksTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS tasks;`)
	return err
}
