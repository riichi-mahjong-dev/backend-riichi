package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSettingsTable, downCreateSettingsTable)
}

func upCreateSettingsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE settings (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255),
				value VARCHAR(255),
				type VARCHAR(255),
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			);
		`)
	return err
}

func downCreateSettingsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE settings;`)
	return err
}
