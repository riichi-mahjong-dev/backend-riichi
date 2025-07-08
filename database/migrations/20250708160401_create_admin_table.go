package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateAdminTable, downCreateAdminTable)
}

func upCreateAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE admins (
				id BIGINT NOT NULL AUTOINCREMENT,
				username VARCHAR(255) NOT NULL,
				password VARCHAR(36) NOT NULL,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
			);
		`)
	return err
}

func downCreateAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE admins;`)
	return err
}
