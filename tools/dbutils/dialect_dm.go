package dbutils

import "fmt"

// DmDialect is a starter implementation for Dameng/DM SQL dialect support.
//
// The SQL snippets below are intentionally conservative and may require
// adjustments depending on the exact DM version and grants.
type DmDialect struct{}

func init() {
	RegisterDMDialect("dm")
}

// RegisterDMDialect registers the dm dialect for the provided driver names.
//
// If no names are provided, it registers for "dm".
func RegisterDMDialect(driverNames ...string) {
	if len(driverNames) == 0 {
		driverNames = []string{"dm"}
	}

	for _, name := range driverNames {
		RegisterDialect(name, DmDialect{})
	}
}

func (d DmDialect) Name() string {
	return "dm"
}

func (d DmDialect) TableColumnsSQL() string {
	return `
SELECT COLUMN_NAME AS name
FROM ALL_TAB_COLUMNS
WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
  AND UPPER(TABLE_NAME) = UPPER({:tableName})
ORDER BY COLUMN_ID`
}

func (d DmDialect) TableInfoSQL() string {
	return `
SELECT
	C.COLUMN_ID - 1 AS cid,
	C.COLUMN_NAME AS name,
	C.DATA_TYPE AS type,
	CASE WHEN C.NULLABLE = 'N' THEN 1 ELSE 0 END AS notnull,
	C.DATA_DEFAULT AS dflt_value,
	CASE
		WHEN EXISTS (
			SELECT 1
			FROM ALL_CONSTRAINTS UC
			JOIN ALL_CONS_COLUMNS UCC ON UC.CONSTRAINT_NAME = UCC.CONSTRAINT_NAME
			WHERE UC.OWNER = C.OWNER
			  AND UC.CONSTRAINT_TYPE = 'P'
			  AND UCC.TABLE_NAME = C.TABLE_NAME
			  AND UCC.COLUMN_NAME = C.COLUMN_NAME
		) THEN 1
		ELSE 0
	END AS pk
FROM ALL_TAB_COLUMNS C
WHERE C.OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
  AND UPPER(C.TABLE_NAME) = UPPER({:tableName})
ORDER BY C.COLUMN_ID`
}

func (d DmDialect) MasterTableName() string {
	// Provides a sqlite_master-like projection (name/sql/type/tbl_name).
	return `
(
	SELECT
		INDEX_NAME AS name,
		NULL AS sql,
		'index' AS type,
		TABLE_NAME AS tbl_name
	FROM ALL_INDEXES
	WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
	UNION ALL
	SELECT
		VIEW_NAME AS name,
		TEXT AS sql,
		'view' AS type,
		VIEW_NAME AS tbl_name
	FROM ALL_VIEWS
	WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
	UNION ALL
	SELECT
		TABLE_NAME AS name,
		NULL AS sql,
		'table' AS type,
		TABLE_NAME AS tbl_name
	FROM ALL_TABLES
	WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
)`
}

func (d DmDialect) SchemaTableName() string {
	// Provides a sqlite_schema-like projection (name/type).
	return `
(
	SELECT
		TABLE_NAME AS name,
		'table' AS type
	FROM ALL_TABLES
	WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
	UNION ALL
	SELECT
		VIEW_NAME AS name,
		'view' AS type
	FROM ALL_VIEWS
	WHERE OWNER = SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA')
)`
}

func (d DmDialect) OptimizeSQL() string {
	// No-op by default. Fill with a DM equivalent if needed.
	return ""
}

func (d DmDialect) WalCheckpointSQL() string {
	// SQLite-specific; no DM equivalent here.
	return ""
}

func (d DmDialect) DateHourExpr(column string) string {
	// The log timestamps are stored as app-formatted strings by default.
	// Keep the grouping expression aligned with the idx_logs_created_hour index.
	return fmt.Sprintf("SUBSTR(%s, 1, 13) || ':00:00'", column)
}

func (d DmDialect) JSONEach(column string) string {
	// TODO: replace with DM JSON table expansion syntax (JSON_TABLE or equivalent).
	return fmt.Sprintf(
		`json_each(CASE WHEN iif(json_valid([[%s]]), json_type([[%s]])='array', FALSE) THEN [[%s]] ELSE json_array([[%s]]) END)`,
		column, column, column, column,
	)
}

func (d DmDialect) JSONArrayLength(column string) string {
	// TODO: replace with DM JSON array length syntax.
	return fmt.Sprintf(
		`json_array_length(CASE WHEN iif(json_valid([[%s]]), json_type([[%s]])='array', FALSE) THEN [[%s]] ELSE (CASE WHEN [[%s]] = '' OR [[%s]] IS NULL THEN json_array() ELSE json_array([[%s]]) END) END)`,
		column, column, column, column, column, column,
	)
}

func (d DmDialect) JSONExtract(column string, path string) string {
	// TODO: replace with DM JSON_VALUE/JSON_QUERY syntax.
	return sqliteDialect{}.JSONExtract(column, path)
}
