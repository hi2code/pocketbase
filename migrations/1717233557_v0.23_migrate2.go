package migrations

import (
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

// note: this migration will be deleted in future version

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		query := "CREATE INDEX IF NOT EXISTS idx__collections_type on {{_collections}} ([[type]]);"
		if dbutils.IsMySQLLike() {
			query = "CREATE INDEX idx__collections_type on {{_collections}} ([[type]]);"
		}

		_, err := txApp.DB().NewQuery(query).Execute()
		if err != nil {
			if dbutils.IsMySQLLike() && strings.Contains(strings.ToLower(err.Error()), "duplicate key name") {
				err = nil
			}
		}
		if err != nil {
			return err
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
