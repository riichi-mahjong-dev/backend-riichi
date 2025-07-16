package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreatePostTable, downCreatePostTable)
}

func upCreatePostTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE posts (
				id BIGINT NOT NULL AUTO_INCREMENT,
				title VARCHAR(255) NOT NULL,
				content TEXT NOT NULL,
				created_by BIGINT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				FOREIGN KEY (created_by) REFERENCES admins(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
	return err
}

func downCreatePostTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE posts;`)
	return err
}
