package dbutils

import "fmt"

// MysqlDialect contains the SQL snippets that differ from the default
// SQLite dialect when PocketBase runs on a MySQL-compatible database.
type MysqlDialect struct{}

func init() {
	RegisterDialect("mysql", MysqlDialect{})
}

func (d MysqlDialect) Name() string {
	return "mysql"
}

func (d MysqlDialect) TableColumnsSQL() string {
	return `
SELECT COLUMN_NAME AS name
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = {:tableName}
ORDER BY ORDINAL_POSITION`
}

func (d MysqlDialect) TableInfoSQL() string {
	return `
SELECT
	ORDINAL_POSITION - 1 AS cid,
	COLUMN_NAME AS name,
	COLUMN_TYPE AS type,
	CASE WHEN IS_NULLABLE = 'NO' THEN 1 ELSE 0 END AS notnull,
	COLUMN_DEFAULT AS dflt_value,
	CASE WHEN COLUMN_KEY = 'PRI' THEN 1 ELSE 0 END AS pk
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = {:tableName}
ORDER BY ORDINAL_POSITION`
}

func (d MysqlDialect) MasterTableName() string {
	return `
(
	SELECT
		INDEX_NAME AS name,
		NULL AS sql,
		'index' AS type,
		TABLE_NAME AS tbl_name
	FROM information_schema.STATISTICS
	WHERE TABLE_SCHEMA = DATABASE()
	UNION ALL
	SELECT
		TABLE_NAME AS name,
		VIEW_DEFINITION AS sql,
		'view' AS type,
		TABLE_NAME AS tbl_name
	FROM information_schema.VIEWS
	WHERE TABLE_SCHEMA = DATABASE()
	UNION ALL
	SELECT
		TABLE_NAME AS name,
		NULL AS sql,
		'table' AS type,
		TABLE_NAME AS tbl_name
	FROM information_schema.TABLES
	WHERE TABLE_SCHEMA = DATABASE()
	  AND TABLE_TYPE = 'BASE TABLE'
)`
}

func (d MysqlDialect) SchemaTableName() string {
	return `
(
	SELECT
		TABLE_NAME AS name,
		CASE
			WHEN TABLE_TYPE = 'BASE TABLE' THEN 'table'
			WHEN TABLE_TYPE = 'VIEW' THEN 'view'
			ELSE LOWER(TABLE_TYPE)
		END AS type
	FROM information_schema.TABLES
	WHERE TABLE_SCHEMA = DATABASE()
)`
}

func (d MysqlDialect) OptimizeSQL() string {
	return ""
}

func (d MysqlDialect) WalCheckpointSQL() string {
	return ""
}

func (d MysqlDialect) DateHourExpr(column string) string {
	return fmt.Sprintf("CONCAT(SUBSTR(%s, 1, 13), ':00:00')", column)
}

func (d MysqlDialect) JSONEach(column string) string {
	// TODO: replace with MySQL JSON_TABLE syntax.
	return sqliteDialect{}.JSONEach(column)
}

func (d MysqlDialect) JSONArrayLength(column string) string {
	// TODO: replace with MySQL JSON_LENGTH syntax.
	return sqliteDialect{}.JSONArrayLength(column)
}

func (d MysqlDialect) JSONExtract(column string, path string) string {
	// TODO: replace with MySQL JSON_EXTRACT/JSON_UNQUOTE syntax.
	return sqliteDialect{}.JSONExtract(column, path)
}
