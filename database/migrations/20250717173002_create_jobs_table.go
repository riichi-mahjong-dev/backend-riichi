package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateJobsTable, downCreateJobsTable)
}

func upCreateJobsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
			CREATE TABLE jobs (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				job_type VARCHAR(50),
				payload JSON,
				status ENUM('queued', 'processing', 'done', 'error') DEFAULT 'queued',
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			);
		`)
	return err
}

func downCreateJobsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE jobs;`)
	return err
}
