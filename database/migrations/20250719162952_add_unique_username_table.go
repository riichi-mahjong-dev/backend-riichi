package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddUniqueUsernameTable, downAddUniqueUsernameTable)
}

func upAddUniqueUsernameTable(tx *sql.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE admins ADD CONSTRAINT unique_username UNIQUE (username);`); err != nil {
		return err
	}

	// Make created_by nullable
	if _, err := tx.Exec(`ALTER TABLE players ADD CONSTRAINT unique_username UNIQUE (username);`); err != nil {
		return err
	}

	return nil
}

func downAddUniqueUsernameTable(tx *sql.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE admins DROP INDEX unique_username;`); err != nil {
		return err
	}

	// Make created_by nullable
	if _, err := tx.Exec(`ALTER TABLE players DROP INDEX unique_username;`); err != nil {
		return err
	}

	return nil
}
