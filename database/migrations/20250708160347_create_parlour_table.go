package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateParlourTable, downCreateParlourTable)
}

func upCreateParlourTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE parlours (
				id BITINT NOT NULL AUTOINCREMENT,
				name VARCHAR(255) NOT NULL,
				country VARCHAR(255) NOT NULL,
				province_id BIGINT NOT NULL,
				address TEXT,
			);
		`)
	return err
}

func downCreateParlourTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE parlours;`)
	return err
}
