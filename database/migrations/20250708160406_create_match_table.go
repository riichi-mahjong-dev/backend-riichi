package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateMatchTable, downCreateMatchTable)
}

func upCreateMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE matches (
				id BIGINT NOT NULL AUTOINCREMENT,
				player_1_id BIGINT,
				player_2_id BIGINT,
				player_3_id BIGINT,
				player_4_id BIGINT,
				player_1_score INT,
				player_2_score INT,
				player_3_score INT,
				player_4_score INT,
				parlour_id BIGINT NOT NULL,
				created_by BIGINT NOT NULL,
				approved_by BIGINT,
				approved_at TIMESTAMP,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
			);
		`)
	return err
}

func downCreateMatchTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE matches`)
	return err
}
