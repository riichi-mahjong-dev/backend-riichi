package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateProvinceTable, downCreateProvinceTable)
}

func upCreateProvinceTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE provinces (
				id BIGINT NOT NULL AUTOINCREMENT,
				name VARCHAR(255) NOT NULL,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
			);
		`)
	return err
}

func downCreateProvinceTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE provinces;`)
	return err
}
