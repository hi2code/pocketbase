package dbutils

import (
	"fmt"
	"strings"
)

type mysqlDialect struct{}

func (d mysqlDialect) Name() string {
	return "mysql"
}

func (d mysqlDialect) TableColumnsSQL() string {
	return `
SELECT COLUMN_NAME AS name
FROM information_schema.columns
WHERE table_schema = DATABASE()
  AND table_name = {:tableName}
ORDER BY ordinal_position`
}

func (d mysqlDialect) TableInfoSQL() string {
	return `
SELECT
	c.ordinal_position - 1 AS cid,
	c.column_name AS name,
	c.column_type AS type,
	CASE WHEN c.is_nullable = 'NO' THEN 1 ELSE 0 END AS notnull,
	c.column_default AS dflt_value,
	CASE WHEN c.column_key = 'PRI' THEN 1 ELSE 0 END AS pk
FROM information_schema.columns c
WHERE c.table_schema = DATABASE()
  AND c.table_name = {:tableName}
ORDER BY c.ordinal_position`
}

func (d mysqlDialect) MasterTableName() string {
	return `
(
	SELECT
			s.index_name AS name,
			CONCAT(
				'CREATE ',
				CASE WHEN s.non_unique = 0 THEN 'UNIQUE ' ELSE '' END,
				'INDEX ', s.index_name, ' ON ', s.table_name, ' (',
				GROUP_CONCAT(s.column_name ORDER BY s.seq_in_index SEPARATOR ', '),
				')'
			) AS sql,
		'index' AS type,
		s.table_name AS tbl_name
	FROM information_schema.statistics s
	WHERE s.table_schema = DATABASE()
	GROUP BY s.table_name, s.index_name, s.non_unique
	UNION ALL
	SELECT
		t.table_name AS name,
		NULL AS sql,
		LOWER(t.table_type) AS type,
		t.table_name AS tbl_name
	FROM information_schema.tables t
	WHERE t.table_schema = DATABASE()
)`
}

func (d mysqlDialect) SchemaTableName() string {
	return `
(
	SELECT
		table_name AS name,
		CASE WHEN table_type = 'VIEW' THEN 'view' ELSE 'table' END AS type
	FROM information_schema.tables
	WHERE table_schema = DATABASE()
)`
}

func (d mysqlDialect) OptimizeSQL() string {
	return ""
}

func (d mysqlDialect) WalCheckpointSQL() string {
	return ""
}

func (d mysqlDialect) DateHourExpr(column string) string {
	return fmt.Sprintf("DATE_FORMAT(%s, '%%Y-%%m-%%d %%H:00:00')", column)
}

func (d mysqlDialect) JSONEach(column string) string {
	normalized := fmt.Sprintf(
		"CASE WHEN JSON_VALID([[%s]]) AND JSON_TYPE([[%s]])='ARRAY' THEN [[%s]] ELSE (CASE WHEN [[%s]] IS NULL OR [[%s]]='' THEN JSON_ARRAY() ELSE JSON_ARRAY([[%s]]) END) END",
		column, column, column, column, column, column,
	)
	return fmt.Sprintf(
		"(SELECT JSON_UNQUOTE(JSON_EXTRACT(%s, CONCAT('$[', _pb_seq.n, ']'))) AS value FROM _pb_seq WHERE _pb_seq.n < JSON_LENGTH(%s))",
		normalized, normalized,
	)
}

func (d mysqlDialect) JSONArrayLength(column string) string {
	return fmt.Sprintf(
		"JSON_LENGTH(CASE WHEN JSON_VALID([[%s]]) AND JSON_TYPE([[%s]])='ARRAY' THEN [[%s]] ELSE (CASE WHEN [[%s]]='' OR [[%s]] IS NULL THEN JSON_ARRAY() ELSE JSON_ARRAY([[%s]]) END) END)",
		column, column, column, column, column, column,
	)
}

func (d mysqlDialect) JSONExtract(column string, path string) string {
	if path != "" && !strings.HasPrefix(path, "[") {
		path = "." + path
	}

	return fmt.Sprintf(
		"(CASE WHEN JSON_VALID([[%s]]) THEN JSON_EXTRACT([[%s]], '$%s') ELSE JSON_EXTRACT(JSON_OBJECT('pb', [[%s]]), '$.pb%s') END)",
		column, column, path, column, path,
	)
}
