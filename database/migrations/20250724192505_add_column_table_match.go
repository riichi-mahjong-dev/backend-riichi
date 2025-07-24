package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddColumnTableMatch, downAddColumnTableMatch)
}

func upAddColumnTableMatch(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE matches ADD COLUMN playing_at TIMESTAMP;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE posts ADD COLUMN is_published TINYINT DEFAULT 0;
	`)
	return err
}

func downAddColumnTableMatch(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE matches DROP COLUMN playing_at;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE posts DROP COLUMN is_published;
	`)
	return err
}
