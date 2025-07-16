package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateProvinceTable, downCreateProvinceTable)
}

func upCreateProvinceTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE provinces (
				id BIGINT NOT NULL AUTO_INCREMENT,
				name VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id)
			);
		`)
	return err
}

func downCreateProvinceTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE provinces;`)
	return err
}
