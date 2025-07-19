package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddColumnCreatedByAdminMatchTable, downAddColumnCreatedByAdminMatchTable)
}

func upAddColumnCreatedByAdminMatchTable(tx *sql.Tx) error {
	// Drop existing FK on created_by
	if _, err := tx.Exec(`ALTER TABLE matches DROP FOREIGN KEY matches_ibfk_2;`); err != nil {
		return err
	}

	// Make created_by nullable
	if _, err := tx.Exec(`ALTER TABLE matches MODIFY COLUMN created_by BIGINT NULL;`); err != nil {
		return err
	}

	// Re-add FK on created_by
	if _, err := tx.Exec(`
		ALTER TABLE matches ADD CONSTRAINT matches_ibfk_2 
		FOREIGN KEY (created_by) REFERENCES players(id)
		ON DELETE CASCADE ON UPDATE CASCADE;
	`); err != nil {
		return err
	}

	return nil
}

func downAddColumnCreatedByAdminMatchTable(tx *sql.Tx) error {
	// Drop FK on created_by_admin
	if _, err := tx.Exec(`
		ALTER TABLE matches DROP FOREIGN KEY matches_ibfk_2;
	`); err != nil {
		return err
	}

	// Drop created_by_admin column
	if _, err := tx.Exec(`
		ALTER TABLE matches MODIFY COLUMN created_by BIGINT NOT NULL;
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		ALTER TABLE matches ADD CONSTRAINT matches_ibfk_2 
		FOREIGN KEY (created_by) REFERENCES players(id)
		ON DELETE CASCADE ON UPDATE CASCADE;
	`); err != nil {
		return err
	}

	return nil
}
