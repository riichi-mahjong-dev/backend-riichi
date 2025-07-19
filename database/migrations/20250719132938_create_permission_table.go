package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreatePermissionTable, downCreatePermissionTable)
}

func upCreatePermissionTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE admin_permissions(
			admin_id BIGINT NOT NULL,
			province_id BIGINT NOT NULL,
			parlour_id BIGINT NOT NULL,
			FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (parlour_id) REFERENCES parlours(id) ON DELETE CASCADE ON UPDATE CASCADE
		);
	`)
	return err
}

func downCreatePermissionTable(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE admin_permissions;`)
	return err
}
