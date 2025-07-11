package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddRoleToAdminTable, downAddRoleToAdminTable)
}

func upAddRoleToAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE admins ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'staff';
	`)
	return err
}

func downAddRoleToAdminTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE admins DROP COLUMN role;
	`)
	return err
}
