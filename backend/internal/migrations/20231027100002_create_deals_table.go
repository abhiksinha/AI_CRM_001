package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDealsTable, downDealsTable)
}

func upDealsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS deals (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			stage VARCHAR(100) NOT NULL,
			value DECIMAL,
			expected_close_date DATE,
			contact_id VARCHAR(36) REFERENCES contacts(id) ON DELETE CASCADE,
			owner_id VARCHAR(36) REFERENCES users(id) ON DELETE SET NULL,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`)
	return err
}

func downDealsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS deals;`)
	return err
}
