package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upFixPasswordColumnSize, downFixPasswordColumnSize)
}

func upFixPasswordColumnSize(tx *sql.Tx) error {
	// Fix admin password column size
	_, err := tx.Exec(`
		ALTER TABLE admins MODIFY COLUMN password VARCHAR(255) NOT NULL;
	`)
	if err != nil {
		return err
	}

	// Fix player password column size
	_, err = tx.Exec(`
		ALTER TABLE players MODIFY COLUMN password VARCHAR(255) NOT NULL;
	`)
	return err
}

func downFixPasswordColumnSize(tx *sql.Tx) error {
	// Revert admin password column size
	_, err := tx.Exec(`
		ALTER TABLE admins MODIFY COLUMN password VARCHAR(36) NOT NULL;
	`)
	if err != nil {
		return err
	}

	// Revert player password column size
	_, err = tx.Exec(`
		ALTER TABLE players MODIFY COLUMN password VARCHAR(36) NOT NULL;
	`)
	return err
}
