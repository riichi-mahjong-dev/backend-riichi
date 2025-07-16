package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateMatchPlayerTable, downCreateMatchPlayerTable)
}

func upCreateMatchPlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE match_players (
				id BIGINT NOT NULL AUTO_INCREMENT,
				match_id BIGINT NOT NULL,
				player_id BIGINT NOT NULL,
				point INT,
				mmr_delta INT,
				pinalty INT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE ON UPDATE CASCADE,
				FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
	return err
}

func downCreateMatchPlayerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE match_players;`)
	return err
}
