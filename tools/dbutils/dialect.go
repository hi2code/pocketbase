package dbutils

import (
	"fmt"
	"strings"
	"sync"
)

// Dialect defines the SQL expressions that differ between database engines.
type Dialect interface {
	Name() string

	// metadata expressions
	TableColumnsSQL() string
	TableInfoSQL() string
	MasterTableName() string
	SchemaTableName() string

	// maintenance expressions
	OptimizeSQL() string
	WalCheckpointSQL() string

	// date/time helpers
	DateHourExpr(column string) string

	// json helpers
	JSONEach(column string) string
	JSONArrayLength(column string) string
	JSONExtract(column string, path string) string
}

var (
	dialectsMu sync.RWMutex
	dialects           = map[string]Dialect{}
	active     Dialect = sqliteDialect{}
)

func init() {
	RegisterDialect("sqlite", sqliteDialect{})
	RegisterDialect("sqlite3", sqliteDialect{})
}

// RegisterDialect registers or overrides a SQL dialect for a driver name.
func RegisterDialect(driver string, dialect Dialect) {
	driver = strings.ToLower(strings.TrimSpace(driver))
	if driver == "" || dialect == nil {
		return
	}

	dialectsMu.Lock()
	defer dialectsMu.Unlock()

	dialects[driver] = dialect
}

// GetDialect returns the currently active dialect.
func GetDialect() Dialect {
	dialectsMu.RLock()
	defer dialectsMu.RUnlock()

	return active
}

// SetDialectByDriver sets the active dialect by the provided driver name.
//
// If no dialect is registered for driver, the default sqlite dialect is used.
func SetDialectByDriver(driver string) {
	driver = strings.ToLower(strings.TrimSpace(driver))

	dialectsMu.Lock()
	defer dialectsMu.Unlock()

	if dialect, ok := dialects[driver]; ok {
		active = dialect
		return
	}

	active = sqliteDialect{}
}

type sqliteDialect struct{}

func (d sqliteDialect) Name() string {
	return "sqlite"
}

func (d sqliteDialect) TableColumnsSQL() string {
	return "SELECT name FROM PRAGMA_TABLE_INFO({:tableName})"
}

func (d sqliteDialect) TableInfoSQL() string {
	return "SELECT * FROM PRAGMA_TABLE_INFO({:tableName})"
}

func (d sqliteDialect) MasterTableName() string {
	return "sqlite_master"
}

func (d sqliteDialect) SchemaTableName() string {
	return "sqlite_schema"
}

func (d sqliteDialect) OptimizeSQL() string {
	return "PRAGMA optimize"
}

func (d sqliteDialect) WalCheckpointSQL() string {
	return "PRAGMA wal_checkpoint(TRUNCATE)"
}

func (d sqliteDialect) DateHourExpr(column string) string {
	return fmt.Sprintf("strftime('%%Y-%%m-%%d %%H:00:00', %s)", column)
}

func (d sqliteDialect) JSONEach(column string) string {
	// note: we are not using the new and shorter "if(x,y)" syntax for
	// compatibility with custom drivers that use older SQLite version
	return fmt.Sprintf(
		`json_each(CASE WHEN iif(json_valid([[%s]]), json_type([[%s]])='array', FALSE) THEN [[%s]] ELSE json_array([[%s]]) END)`,
		column, column, column, column,
	)
}

func (d sqliteDialect) JSONArrayLength(column string) string {
	// note: we are not using the new and shorter "if(x,y)" syntax for
	// compatibility with custom drivers that use older SQLite version
	return fmt.Sprintf(
		`json_array_length(CASE WHEN iif(json_valid([[%s]]), json_type([[%s]])='array', FALSE) THEN [[%s]] ELSE (CASE WHEN [[%s]] = '' OR [[%s]] IS NULL THEN json_array() ELSE json_array([[%s]]) END) END)`,
		column, column, column, column, column, column,
	)
}

func (d sqliteDialect) JSONExtract(column string, path string) string {
	// prefix the path with dot if it is not starting with array notation
	if path != "" && !strings.HasPrefix(path, "[") {
		path = "." + path
	}

	return fmt.Sprintf(
		// note: the extra object wrapping is needed to workaround the cases where a json_extract is used with non-json columns.
		"(CASE WHEN json_valid([[%s]]) THEN JSON_EXTRACT([[%s]], '$%s') ELSE JSON_EXTRACT(json_object('pb', [[%s]]), '$.pb%s') END)",
		column,
		column,
		path,
		column,
		path,
	)
}
