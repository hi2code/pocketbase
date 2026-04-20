package core

import (
	"database/sql"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
)

// TableColumns returns all column names of a single table by its name.
func (app *BaseApp) TableColumns(tableName string) ([]string, error) {
	columns := []string{}

	err := app.ConcurrentDB().NewQuery(dbutils.GetDialect().TableColumnsSQL()).
		Bind(dbx.Params{"tableName": tableName}).
		Column(&columns)

	return columns, err
}

type TableInfoRow struct {
	// the `db:"pk"` tag has special semantic so we cannot rename
	// the original field without specifying a custom mapper
	PK int

	Index        int            `db:"cid"`
	Name         string         `db:"name"`
	Type         string         `db:"type"`
	NotNull      bool           `db:"notnull"`
	DefaultValue sql.NullString `db:"dflt_value"`
}

// TableInfo returns the "table_info" pragma result for the specified table.
func (app *BaseApp) TableInfo(tableName string) ([]*TableInfoRow, error) {
	info := []*TableInfoRow{}

	err := app.ConcurrentDB().NewQuery(dbutils.GetDialect().TableInfoSQL()).
		Bind(dbx.Params{"tableName": tableName}).
		All(&info)
	if err != nil {
		return nil, err
	}

	// mattn/go-sqlite3 doesn't throw an error on invalid or missing table
	// so we additionally have to check whether the loaded info result is nonempty
	if len(info) == 0 {
		return nil, fmt.Errorf("empty table info probably due to invalid or missing table %s", tableName)
	}

	return info, nil
}

// TableIndexes returns a name grouped map with all non empty index of the specified table.
//
// Note: This method doesn't return an error on nonexisting table.
func (app *BaseApp) TableIndexes(tableName string) (map[string]string, error) {
	indexes := []struct {
		Name string `db:"name"`
		Sql  string `db:"sql"`
	}{}

	dialectName := dbutils.GetDialect().Name()

	var err error
	switch dialectName {
	case "mysql", "dm":
		mysqlIndexes := []struct {
			Name string `db:"name"`
			Sql  string `db:"index_sql"`
		}{}

		err = app.ConcurrentDB().NewQuery(`
			SELECT index_name AS name, index_name AS index_sql
			FROM information_schema.statistics
			WHERE table_schema = DATABASE()
			  AND table_name = {:tableName}
			  AND index_name <> 'PRIMARY'
			GROUP BY index_name
		`).
			Bind(dbx.Params{"tableName": tableName}).
			All(&mysqlIndexes)
		if err == nil {
			indexes = make([]struct {
				Name string `db:"name"`
				Sql  string `db:"sql"`
			}, len(mysqlIndexes))
			for i, idx := range mysqlIndexes {
				indexes[i] = struct {
					Name string `db:"name"`
					Sql  string `db:"sql"`
				}{
					Name: idx.Name,
					Sql:  idx.Sql,
				}
			}
		}
	case "pg":
		err = app.ConcurrentDB().NewQuery(`
			SELECT indexname AS name, indexdef AS sql
			FROM pg_indexes
			WHERE schemaname = current_schema()
			  AND tablename = {:tableName}
		`).
			Bind(dbx.Params{"tableName": tableName}).
			All(&indexes)
	default:
		err = app.ConcurrentDB().Select("name", "sql").
			From(dbutils.GetDialect().MasterTableName()).
			AndWhere(dbx.NewExp("sql is not null")).
			AndWhere(dbx.HashExp{
				"type":     "index",
				"tbl_name": tableName,
			}).
			All(&indexes)
	}
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(indexes))

	for _, idx := range indexes {
		result[idx.Name] = idx.Sql
	}

	return result, nil
}

// DeleteTable drops the specified table.
//
// This method is a no-op if a table with the provided name doesn't exist.
//
// NB! Be aware that this method is vulnerable to SQL injection and the
// "dangerousTableName" argument must come only from trusted input!
func (app *BaseApp) DeleteTable(dangerousTableName string) error {
	_, err := app.NonconcurrentDB().NewQuery(fmt.Sprintf(
		"DROP TABLE IF EXISTS {{%s}}",
		dangerousTableName,
	)).Execute()

	return err
}

// HasTable checks if a table (or view) with the provided name exists (case insensitive).
// in the data.db.
func (app *BaseApp) HasTable(tableName string) bool {
	return app.hasTable(app.ConcurrentDB(), tableName)
}

// AuxHasTable checks if a table (or view) with the provided name exists (case insensitive)
// in the auixiliary.db.
func (app *BaseApp) AuxHasTable(tableName string) bool {
	return app.hasTable(app.AuxConcurrentDB(), tableName)
}

func (app *BaseApp) hasTable(db dbx.Builder, tableName string) bool {
	var exists int

	dialectName := dbutils.GetDialect().Name()

	var err error

	switch dialectName {
	case "mysql", "dm":
		err = db.NewQuery(`
			SELECT (1)
			FROM information_schema.tables
			WHERE table_schema = DATABASE()
			  AND table_type IN ('BASE TABLE', 'VIEW')
			  AND LOWER(table_name) = LOWER({:tableName})
			LIMIT 1
		`).
			Bind(dbx.Params{"tableName": tableName}).
			Row(&exists)
	case "pg":
		err = db.NewQuery(`
			SELECT (1)
			FROM information_schema.tables
			WHERE table_schema = current_schema()
			  AND table_type IN ('BASE TABLE', 'VIEW')
			  AND LOWER(table_name) = LOWER({:tableName})
			LIMIT 1
		`).
			Bind(dbx.Params{"tableName": tableName}).
			Row(&exists)
	default:
		err = db.Select("(1)").
			From(dbutils.GetDialect().SchemaTableName()).
			AndWhere(dbx.HashExp{"type": []any{"table", "view"}}).
			AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
			Limit(1).
			Row(&exists)
	}

	if err != nil {
		app.Logger().Debug("hasTable probe failed", "dialect", dialectName, "table", tableName, "error", err)
		return false
	}

	return exists > 0
}

// Vacuum executes VACUUM on the data.db in order to reclaim unused data db disk space.
func (app *BaseApp) Vacuum() error {
	return app.vacuum(app.NonconcurrentDB())
}

// AuxVacuum executes VACUUM on the auxiliary.db in order to reclaim unused auxiliary db disk space.
func (app *BaseApp) AuxVacuum() error {
	return app.vacuum(app.AuxNonconcurrentDB())
}

func (app *BaseApp) vacuum(db dbx.Builder) error {
	if dbutils.GetDialect().Name() != "sqlite" {
		return nil
	}

	_, err := db.NewQuery("VACUUM").Execute()

	return err
}
