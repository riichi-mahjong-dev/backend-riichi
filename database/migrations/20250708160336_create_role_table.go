package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateRoleTable, downCreateRoleTable)
}

func upCreateRoleTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE roles (
				id BIGINT NOT NULL AUTO_INCREMENT,
				name VARCHAR(255) NOT NULL,
				guard_name VARCHAR(255) NOT NULL,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
				PRIMARY KEY (id)
			);
		`)
	return err
}

func downCreateRoleTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE roles;`)
	return err
}
