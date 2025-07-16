package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateMatchTable, downCreateMatchTable)
}

func upCreateMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE matches (
				id BIGINT NOT NULL AUTO_INCREMENT,
				parlour_id BIGINT NOT NULL,
				created_by BIGINT NOT NULL,
				approved_by BIGINT,
				approved_at TIMESTAMP,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				FOREIGN KEY (parlour_id) REFERENCES parlours(id) ON DELETE CASCADE ON UPDATE CASCADE,
				FOREIGN KEY (created_by) REFERENCES players(id) ON DELETE CASCADE ON UPDATE CASCADE,
				FOREIGN KEY (approved_by) REFERENCES admins(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
	return err
}

func downCreateMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE matches;`)
	return err
}
