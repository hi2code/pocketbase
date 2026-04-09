package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func init() {
	core.SystemMigrations.Add(&core.Migration{
		Up: func(txApp core.App) error {
			if dbutils.GetDialect().Name() == "dm" {
				_, execErr := txApp.AuxDB().NewQuery(`
					CREATE TABLE {{_logs}} (
						[[id]]      VARCHAR(32) PRIMARY KEY NOT NULL,
						[[level]]   INTEGER DEFAULT 0 NOT NULL,
						[[message]] VARCHAR(4000) DEFAULT '' NOT NULL,
						[[data]]    VARCHAR(4000) DEFAULT '{}' NOT NULL,
						[[created]] VARCHAR(32) DEFAULT '' NOT NULL
					);

					CREATE INDEX idx_logs_level ON {{_logs}} ([[level]]);
					CREATE INDEX idx_logs_message ON {{_logs}} ([[message]]);
					CREATE INDEX idx_logs_created_hour ON {{_logs}} ((SUBSTR([[created]], 1, 13) || ':00:00'));
				`).Execute()
				return execErr
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
