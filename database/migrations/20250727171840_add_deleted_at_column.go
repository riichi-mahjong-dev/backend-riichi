package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddDeletedAtColumn, downAddDeletedAtColumn)
}

func upAddDeletedAtColumn(tx *sql.Tx) error {
	queries := []string{
		`ALTER TABLE matches ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_matches_deleted_at ON matches(deleted_at);`,

		`ALTER TABLE players ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_players_deleted_at ON players(deleted_at);`,

		`ALTER TABLE provinces ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_provinces_deleted_at ON provinces(deleted_at);`,

		`ALTER TABLE settings ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_settings_deleted_at ON settings(deleted_at);`,

		`ALTER TABLE posts ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);`,

		`ALTER TABLE admins ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_admins_deleted_at ON admins(deleted_at);`,

		`ALTER TABLE parlours ADD COLUMN deleted_at TIMESTAMP NULL;`,
		`CREATE INDEX idx_parlours_deleted_at ON parlours(deleted_at);`,
	}

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func downAddDeletedAtColumn(tx *sql.Tx) error {
	queries := []string{
		`DROP INDEX idx_matches_deleted_at ON matches;`,
		`ALTER TABLE matches DROP COLUMN deleted_at;`,

		`DROP INDEX idx_players_deleted_at ON players;`,
		`ALTER TABLE players DROP COLUMN deleted_at;`,

		`DROP INDEX idx_provinces_deleted_at ON provinces;`,
		`ALTER TABLE provinces DROP COLUMN deleted_at;`,

		`DROP INDEX idx_settings_deleted_at ON settings;`,
		`ALTER TABLE settings DROP COLUMN deleted_at;`,

		`DROP INDEX idx_posts_deleted_at ON posts;`,
		`ALTER TABLE posts DROP COLUMN deleted_at;`,

		`DROP INDEX idx_admins_deleted_at ON admins;`,
		`ALTER TABLE admins DROP COLUMN deleted_at;`,

		`DROP INDEX idx_parlours_deleted_at ON parlours;`,
		`ALTER TABLE parlours DROP COLUMN deleted_at;`,
	}

	for _, q := range queries {
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
