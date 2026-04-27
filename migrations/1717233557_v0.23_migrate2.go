package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

// note: this migration will be deleted in future version

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		if dbutils.GetDialect().Name() == "dm" {
			err := dmCreateIndexIfMissing(
				txApp.DB(),
				"_collections",
				"idx__collections_type",
				"CREATE INDEX idx__collections_type ON [[_collections]] ([[type]])",
			)
			if err != nil {
				return err
			}
		} else if dbutils.GetDialect().Name() != "mysql" {
			_, err := txApp.DB().NewQuery("CREATE INDEX IF NOT EXISTS idx__collections_type on {{_collections}} ([[type]]);").Execute()
			if err != nil {
				return err
			}
		}

		// reset mfas and otps delete rule
		collectionNames := []string{core.CollectionNameMFAs, core.CollectionNameOTPs}
		for _, name := range collectionNames {
			col, err := txApp.FindCollectionByNameOrId(name)
			if err != nil {
				return err
			}

			if col.DeleteRule != nil {
				col.DeleteRule = nil
				err = txApp.SaveNoValidate(col)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}, nil)
}
