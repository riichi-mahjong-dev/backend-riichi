package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddColumnPostTable, downAddColumnPostTable)
}

func upAddColumnPostTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE posts ADD COLUMN slug VARCHAR(255);
	`)
	return err
}

func downAddColumnPostTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE posts DROP COLUMN slug;
	`)
	return err
}
