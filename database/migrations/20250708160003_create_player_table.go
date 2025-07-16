package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreatePlayerTable, downCreatePlayerTable)
}

func upCreatePlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE players (
				id BIGINT NOT NULL AUTO_INCREMENT,
				username VARCHAR(255) NOT NULL,
				password VARCHAR(36) NOT NULL,
				rank INT NOT NULL,
				country VARCHAR(255) NOT NULL,
				province_id BIGINT NOT NULL,
				name VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
	return err
}

func downCreatePlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE players;`)
	return err
}
