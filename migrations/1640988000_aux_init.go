package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func init() {
	core.SystemMigrations.Add(&core.Migration{
		Up: func(txApp core.App) error {
			if dbutils.GetDialect().Name() == "mysql" {
				var logsTableCount int
				if err := txApp.AuxDB().NewQuery(`
					SELECT COUNT(1)
					FROM information_schema.TABLES
					WHERE TABLE_SCHEMA = DATABASE()
					  AND TABLE_NAME = '_logs'
				`).Row(&logsTableCount); err != nil {
					return err
				}
				if logsTableCount > 0 {
					return nil
				}

				if _, err := txApp.AuxDB().NewQuery(`
					CREATE TABLE IF NOT EXISTS {{_logs}} (
						[[id]]      VARCHAR(32) PRIMARY KEY NOT NULL,
						[[level]]   INTEGER DEFAULT 0 NOT NULL,
						[[message]] TEXT NOT NULL,
						[[data]]    JSON DEFAULT NULL,
						[[created]] VARCHAR(32) DEFAULT '' NOT NULL
					)
				`).Execute(); err != nil {
					return err
				}

				for _, query := range []string{
					`CREATE INDEX idx_logs_level ON {{_logs}} ([[level]])`,
					`CREATE INDEX idx_logs_message ON {{_logs}} ([[message]](255))`,
					`CREATE INDEX idx_logs_created_hour ON {{_logs}} ([[created]])`,
				} {
					if _, err := txApp.AuxDB().NewQuery(query).Execute(); err != nil {
						return err
					}
				}

				return nil
			}

			if dbutils.GetDialect().Name() == "dm" {
				var logsTableCount int
				if err := txApp.AuxDB().NewQuery(`
					SELECT COUNT(1)
					FROM USER_TABLES
					WHERE UPPER(TABLE_NAME) = UPPER('_logs')
				`).Row(&logsTableCount); err != nil {
					return err
				}
				if logsTableCount > 0 {
					return nil
				}

				if _, err := txApp.AuxDB().NewQuery(`
					CREATE TABLE [[_logs]] (
						[[id]]      VARCHAR(32) PRIMARY KEY NOT NULL,
						[[level]]   INTEGER DEFAULT 0 NOT NULL,
						[[message]] VARCHAR(4000) DEFAULT '' NOT NULL,
						[[data]]    VARCHAR(4000) DEFAULT '{}' NOT NULL,
						[[created]] VARCHAR(32) DEFAULT '' NOT NULL
					)
				`).Execute(); err != nil {
					return err
				}

				indexes := []struct {
					name  string
					query string
				}{
					{"idx_logs_level", `CREATE INDEX idx_logs_level ON [[_logs]] ([[level]])`},
					{"idx_logs_message", `CREATE INDEX idx_logs_message ON [[_logs]] ([[message]])`},
					{"idx_logs_created_hour", `CREATE INDEX idx_logs_created_hour ON [[_logs]] ([[created]])`},
				}
				for _, index := range indexes {
					if err := dmCreateIndexIfMissing(txApp.AuxDB(), "_logs", index.name, index.query); err != nil {
						return err
					}
				}

				return nil
			}

			_, execErr := txApp.AuxDB().NewQuery(`
				CREATE TABLE IF NOT EXISTS {{_logs}} (
					[[id]]      TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL,
					[[level]]   INTEGER DEFAULT 0 NOT NULL,
					[[message]] TEXT DEFAULT "" NOT NULL,
					[[data]]    JSON DEFAULT "{}" NOT NULL,
					[[created]] TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
				);

				CREATE INDEX IF NOT EXISTS idx_logs_level on {{_logs}} ([[level]]);
				CREATE INDEX IF NOT EXISTS idx_logs_message on {{_logs}} ([[message]]);
				CREATE INDEX IF NOT EXISTS idx_logs_created_hour on {{_logs}} (strftime('%Y-%m-%d %H:00:00', [[created]]));
			`).Execute()

			return execErr
		},
		Down: func(txApp core.App) error {
			_, err := txApp.AuxDB().DropTable("_logs").Execute()
			return err
		},
		ReapplyCondition: func(txApp core.App, runner *core.MigrationsRunner, fileName string) (bool, error) {
			// reapply only if the _logs table doesn't exist
			exists := txApp.AuxHasTable("_logs")
			return !exists, nil
		},
	})
}
