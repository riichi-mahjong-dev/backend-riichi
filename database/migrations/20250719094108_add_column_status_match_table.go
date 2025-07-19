package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddColumnStatusMatchTable, downAddColumnStatusMatchTable)
}

func upAddColumnStatusMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE matches ADD COLUMN status ENUM('not_sync', 'syncing', 'done') DEFAULT 'not_sync' AFTER approved_at;
	`)
	return err
}

func downAddColumnStatusMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE matches DROP COLUMN status;
	`)
	return err
}
