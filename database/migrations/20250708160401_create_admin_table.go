package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateAdminTable, downCreateAdminTable)
}

func upCreateAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE admins (
				id BIGINT NOT NULL AUTO_INCREMENT,
				username VARCHAR(255) NOT NULL,
				password VARCHAR(36) NOT NULL,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
				PRIMARY KEY (id)
			);
		`)
	return err
}

func downCreateAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE admins;`)
	return err
}
