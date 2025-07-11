package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateParlourTable, downCreateParlourTable)
}

func upCreateParlourTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE parlours (
				id BIGINT NOT NULL AUTO_INCREMENT,
				name VARCHAR(255) NOT NULL,
				country VARCHAR(255) NOT NULL,
				province_id BIGINT NOT NULL,
				address TEXT,
				PRIMARY KEY (id)
			);
		`)
	return err
}

func downCreateParlourTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE parlours;`)
	return err
}
