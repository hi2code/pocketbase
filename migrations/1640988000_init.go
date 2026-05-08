package migrations

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Register is a short alias for `AppMigrations.Register()`
// that is usually used in external/user defined migrations.
func Register(
	up func(app core.App) error,
	down func(app core.App) error,
	optFilename ...string,
) {
	var optFiles []string
	if len(optFilename) > 0 {
		optFiles = optFilename
	} else {
		_, path, _, _ := runtime.Caller(1)
		optFiles = append(optFiles, filepath.Base(path))
	}
	core.AppMigrations.Register(up, down, optFiles...)
}

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		if err := createParamsTable(txApp); err != nil {
			return fmt.Errorf("_params exec error: %w", err)
		}

		// -----------------------------------------------------------

		if dbutils.GetDialect().Name() == "mysql" {
			exists, err := mysqlTableExists(txApp, "_collections")
			if err != nil {
				return fmt.Errorf("_collections exec error: %w", err)
			}
			if !exists {
				if _, execerr := txApp.DB().NewQuery(`
					CREATE TABLE {{_collections}} (
						[[id]]         VARCHAR(32) PRIMARY KEY NOT NULL,
						[[system]]     BOOLEAN DEFAULT FALSE NOT NULL,
						[[type]]       VARCHAR(64) DEFAULT 'base' NOT NULL,
						[[name]]       VARCHAR(255) UNIQUE NOT NULL,
						[[fields]]     JSON DEFAULT NULL,
						[[indexes]]    JSON DEFAULT NULL,
						[[listRule]]   TEXT DEFAULT NULL,
						[[viewRule]]   TEXT DEFAULT NULL,
						[[createRule]] TEXT DEFAULT NULL,
						[[updateRule]] TEXT DEFAULT NULL,
						[[deleteRule]] TEXT DEFAULT NULL,
						[[options]]    JSON DEFAULT NULL,
						[[created]]    VARCHAR(32) DEFAULT '' NOT NULL,
						[[updated]]    VARCHAR(32) DEFAULT '' NOT NULL
					)
				`).Execute(); execerr != nil {
					return fmt.Errorf("_collections exec error: %w", execerr)
				}

				_, execerr := txApp.DB().NewQuery(`CREATE INDEX idx__collections_type on {{_collections}} ([[type]])`).Execute()
				if execerr != nil {
					return fmt.Errorf("_collections exec error: %w", execerr)
				}
			}
		} else if dbutils.GetDialect().Name() == "dm" {
			exists, err := dmTableExists(txApp, "_collections")
			if err != nil {
				return fmt.Errorf("_collections exec error: %w", err)
			}
			if !exists {
				if _, execerr := txApp.DB().NewQuery(`
					CREATE TABLE [[_collections]] (
						[[id]]         VARCHAR(32) PRIMARY KEY NOT NULL,
						[[system]]     SMALLINT DEFAULT 0 NOT NULL,
						[[type]]       VARCHAR(64) DEFAULT 'base' NOT NULL,
						[[name]]       VARCHAR(255) NOT NULL,
						[[fields]]     CLOB,
						[[indexes]]    CLOB,
						[[listRule]]   CLOB,
						[[viewRule]]   CLOB,
						[[createRule]] CLOB,
						[[updateRule]] CLOB,
						[[deleteRule]] CLOB,
						[[options]]    CLOB,
						[[created]]    VARCHAR(32) DEFAULT '' NOT NULL,
						[[updated]]    VARCHAR(32) DEFAULT '' NOT NULL
					)
				`).Execute(); execerr != nil {
					return fmt.Errorf("_collections exec error: %w", execerr)
				}
			}

			indexes := []struct {
				name  string
				query string
			}{
				{"idx__collections_name", `CREATE UNIQUE INDEX idx__collections_name ON [[_collections]] ([[name]])`},
				{"idx__collections_type", `CREATE INDEX idx__collections_type ON [[_collections]] ([[type]])`},
			}
			for _, index := range indexes {
				if execerr := dmCreateIndexIfMissing(txApp.DB(), "_collections", index.name, index.query); execerr != nil {
					return fmt.Errorf("_collections exec error: %w", execerr)
				}
			}
		} else {
			_, execerr := txApp.DB().NewQuery(`
			CREATE TABLE {{_collections}} (
				[[id]]         TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL,
				[[system]]     BOOLEAN DEFAULT FALSE NOT NULL,
				[[type]]       TEXT DEFAULT "base" NOT NULL,
				[[name]]       TEXT UNIQUE NOT NULL,
				[[fields]]     JSON DEFAULT "[]" NOT NULL,
				[[indexes]]    JSON DEFAULT "[]" NOT NULL,
				[[listRule]]   TEXT DEFAULT NULL,
				[[viewRule]]   TEXT DEFAULT NULL,
				[[createRule]] TEXT DEFAULT NULL,
				[[updateRule]] TEXT DEFAULT NULL,
				[[deleteRule]] TEXT DEFAULT NULL,
				[[options]]    JSON DEFAULT "{}" NOT NULL,
				[[created]]    TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				[[updated]]    TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);

			CREATE INDEX IF NOT EXISTS idx__collections_type on {{_collections}} ([[type]]);
		`).Execute()
			if execerr != nil {
				return fmt.Errorf("_collections exec error: %w", execerr)
			}
		}

		if err := createMFAsCollection(txApp); err != nil {
			return fmt.Errorf("_mfas error: %w", err)
		}

		if err := createOTPsCollection(txApp); err != nil {
			return fmt.Errorf("_otps error: %w", err)
		}

		if err := createExternalAuthsCollection(txApp); err != nil {
			return fmt.Errorf("_externalAuths error: %w", err)
		}

		if err := createAuthOriginsCollection(txApp); err != nil {
			return fmt.Errorf("_authOrigins error: %w", err)
		}

		if err := createSuperusersCollection(txApp); err != nil {
			return fmt.Errorf("_superusers error: %w", err)
		}

		if err := createUsersCollection(txApp); err != nil {
			return fmt.Errorf("users error: %w", err)
		}

		return nil
	}, func(txApp core.App) error {
		tables := []string{
			"users",
			core.CollectionNameSuperusers,
			core.CollectionNameMFAs,
			core.CollectionNameOTPs,
			core.CollectionNameAuthOrigins,
			"_params",
			"_collections",
		}

		for _, name := range tables {
			if _, err := txApp.DB().DropTable(name).Execute(); err != nil {
				return err
			}
		}

		return nil
	})
}

func createParamsTable(txApp core.App) error {
	if dbutils.GetDialect().Name() == "mysql" {
		exists, err := mysqlTableExists(txApp, "_params")
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		_, execErr := txApp.DB().NewQuery(`
			CREATE TABLE {{_params}} (
				[[id]]      VARCHAR(32) PRIMARY KEY NOT NULL,
				[[value]]   JSON DEFAULT NULL,
				[[created]] VARCHAR(32) DEFAULT '' NOT NULL,
				[[updated]] VARCHAR(32) DEFAULT '' NOT NULL
			)
		`).Execute()

		return execErr
	}

	if dbutils.GetDialect().Name() == "dm" {
		exists, err := dmTableExists(txApp, "_params")
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		_, execErr := txApp.DB().NewQuery(`
			CREATE TABLE [[_params]] (
				[[id]]      VARCHAR(32) PRIMARY KEY NOT NULL,
				[[value]]   VARCHAR(4000),
				[[created]] VARCHAR(32) DEFAULT '' NOT NULL,
				[[updated]] VARCHAR(32) DEFAULT '' NOT NULL
			)
		`).Execute()

		return execErr
	}

	_, execErr := txApp.DB().NewQuery(`
		CREATE TABLE {{_params}} (
			[[id]]      TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL,
			[[value]]   JSON DEFAULT NULL,
			[[created]] TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
			[[updated]] TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
		);
	`).Execute()

	return execErr
}

func mysqlTableExists(txApp core.App, tableName string) (bool, error) {
	var count int
	err := txApp.DB().NewQuery(`
		SELECT COUNT(1)
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = {:tableName}
	`).Bind(map[string]any{"tableName": tableName}).Row(&count)

	return count > 0, err
}

func dmTableExists(txApp core.App, tableName string) (bool, error) {
	var count int
	err := txApp.DB().NewQuery(`
		SELECT COUNT(1)
		FROM ALL_TABLES
		WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
		  AND UPPER(TABLE_NAME) = UPPER({:tableName})
	`).Bind(map[string]any{"tableName": tableName}).Row(&count)

	return count > 0, err
}

func dmIndexExists(db dbx.Builder, tableName string, indexName string) (bool, error) {
	var count int
	err := db.NewQuery(`
		SELECT COUNT(1)
		FROM ALL_INDEXES
		WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
		  AND UPPER(TABLE_NAME) = UPPER({:tableName})
		  AND UPPER(INDEX_NAME) = UPPER({:indexName})
	`).Bind(dbx.Params{
		"tableName": tableName,
		"indexName": indexName,
	}).Row(&count)

	return count > 0, err
}

func dmCreateIndexIfMissing(db dbx.Builder, tableName string, indexName string, query string) error {
	exists, err := dmIndexExists(db, tableName, indexName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = db.NewQuery(query).Execute()
	return err
}

func createMFAsCollection(txApp core.App) error {
	col := core.NewBaseCollection(core.CollectionNameMFAs)
	col.System = true

	ownerRule := "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
	col.ListRule = types.Pointer(ownerRule)
	col.ViewRule = types.Pointer(ownerRule)

	col.Fields.Add(&core.TextField{
		Name:     "collectionRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "recordRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "method",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	col.AddIndex("idx_mfas_collectionRef_recordRef", false, "collectionRef,recordRef", "")

	return txApp.Save(col)
}

func createOTPsCollection(txApp core.App) error {
	col := core.NewBaseCollection(core.CollectionNameOTPs)
	col.System = true

	ownerRule := "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
	col.ListRule = types.Pointer(ownerRule)
	col.ViewRule = types.Pointer(ownerRule)

	col.Fields.Add(&core.TextField{
		Name:     "collectionRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "recordRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.PasswordField{
		Name:     "password",
		System:   true,
		Hidden:   true,
		Required: true,
		Cost:     8, // low cost for better performance and because it is not critical
	})
	col.Fields.Add(&core.TextField{
		Name:   "sentTo",
		System: true,
		Hidden: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	col.AddIndex("idx_otps_collectionRef_recordRef", false, "collectionRef, recordRef", "")

	return txApp.Save(col)
}

func createAuthOriginsCollection(txApp core.App) error {
	col := core.NewBaseCollection(core.CollectionNameAuthOrigins)
	col.System = true

	ownerRule := "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
	col.ListRule = types.Pointer(ownerRule)
	col.ViewRule = types.Pointer(ownerRule)
	col.DeleteRule = types.Pointer(ownerRule)

	col.Fields.Add(&core.TextField{
		Name:     "collectionRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "recordRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "fingerprint",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	col.AddIndex("idx_authOrigins_unique_pairs", true, "collectionRef, recordRef, fingerprint", "")

	return txApp.Save(col)
}

func createExternalAuthsCollection(txApp core.App) error {
	col := core.NewBaseCollection(core.CollectionNameExternalAuths)
	col.System = true

	ownerRule := "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
	col.ListRule = types.Pointer(ownerRule)
	col.ViewRule = types.Pointer(ownerRule)
	col.DeleteRule = types.Pointer(ownerRule)

	col.Fields.Add(&core.TextField{
		Name:     "collectionRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "recordRef",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "provider",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.TextField{
		Name:     "providerId",
		System:   true,
		Required: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	col.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	col.AddIndex("idx_externalAuths_record_provider", true, "collectionRef, recordRef, provider", "")
	col.AddIndex("idx_externalAuths_collection_provider", true, "collectionRef, provider, providerId", "")

	return txApp.Save(col)
}

func createSuperusersCollection(txApp core.App) error {
	superusers := core.NewAuthCollection(core.CollectionNameSuperusers)
	superusers.System = true
	superusers.Fields.Add(&core.EmailField{
		Name:     "email",
		System:   true,
		Required: true,
	})
	superusers.Fields.Add(&core.AutodateField{
		Name:     "created",
		System:   true,
		OnCreate: true,
	})
	superusers.Fields.Add(&core.AutodateField{
		Name:     "updated",
		System:   true,
		OnCreate: true,
		OnUpdate: true,
	})
	superusers.AuthToken.Duration = 86400 // 1 day

	return txApp.Save(superusers)
}

func createUsersCollection(txApp core.App) error {
	users := core.NewAuthCollection("users", "_pb_users_auth_")

	ownerRule := "id = @request.auth.id"
	users.ListRule = types.Pointer(ownerRule)
	users.ViewRule = types.Pointer(ownerRule)
	users.CreateRule = types.Pointer("")
	users.UpdateRule = types.Pointer(ownerRule)
	users.DeleteRule = types.Pointer(ownerRule)

	users.Fields.Add(&core.TextField{
		Name: "name",
		Max:  255,
	})
	users.Fields.Add(&core.FileField{
		Name:      "avatar",
		MaxSelect: 1,
		MimeTypes: []string{"image/jpeg", "image/png", "image/svg+xml", "image/gif", "image/webp"},
	})
	users.Fields.Add(&core.AutodateField{
		Name:     "created",
		OnCreate: true,
	})
	users.Fields.Add(&core.AutodateField{
		Name:     "updated",
		OnCreate: true,
		OnUpdate: true,
	})
	users.OAuth2.MappedFields.Name = "name"
	users.OAuth2.MappedFields.AvatarURL = "avatar"

	return txApp.Save(users)
}
