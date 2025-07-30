package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddColumnPlayerTable, downAddColumnPlayerTable)
}

func upAddColumnPlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE players ADD COLUMN game_count INT;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE players ADD COLUMN last_login_at DATETIME;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE match_players ADD COLUMN games_played INT;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE players ADD COLUMN last_match_at DATETIME;
	`)
	return err
}

func downAddColumnPlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE players DROP COLUMN game_count;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE players DROP COLUMN last_login_at;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE match_players DROP COLUMN games_played;
	`)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		ALTER TABLE players DROP COLUMN last_match_at;
	`)
	return err
}
