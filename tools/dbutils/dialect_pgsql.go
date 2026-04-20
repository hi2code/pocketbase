package dbutils

import (
	"fmt"
	"strings"
)

type pgsqlDialect struct{}

func (d pgsqlDialect) Name() string {
	return "pg"
}

func (d pgsqlDialect) TableColumnsSQL() string {
	return `
SELECT column_name AS name
FROM information_schema.columns
WHERE table_schema = current_schema()
  AND table_name = {:tableName}
ORDER BY ordinal_position`
}

func (d pgsqlDialect) TableInfoSQL() string {
	return `
SELECT
	c.ordinal_position - 1 AS cid,
	c.column_name AS name,
	c.data_type AS type,
	CASE WHEN c.is_nullable = 'NO' THEN 1 ELSE 0 END AS notnull,
	c.column_default AS dflt_value,
	CASE
		WHEN EXISTS (
			SELECT 1
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
			  ON tc.constraint_name = kcu.constraint_name
			 AND tc.table_schema = kcu.table_schema
			 AND tc.table_name = kcu.table_name
			WHERE tc.constraint_type = 'PRIMARY KEY'
			  AND tc.table_schema = c.table_schema
			  AND tc.table_name = c.table_name
			  AND kcu.column_name = c.column_name
		) THEN 1
		ELSE 0
	END AS pk
FROM information_schema.columns c
WHERE c.table_schema = current_schema()
  AND c.table_name = {:tableName}
ORDER BY c.ordinal_position`
}

func (d pgsqlDialect) MasterTableName() string {
	return `
(
	SELECT
		indexname AS name,
		indexdef AS sql,
		'index' AS type,
		tablename AS tbl_name
	FROM pg_indexes
	WHERE schemaname = current_schema()
	UNION ALL
	SELECT
		table_name AS name,
		NULL AS sql,
		'table' AS type,
		table_name AS tbl_name
	FROM information_schema.tables
	WHERE table_schema = current_schema() AND table_type = 'BASE TABLE'
	UNION ALL
	SELECT
		table_name AS name,
		NULL AS sql,
		'view' AS type,
		table_name AS tbl_name
	FROM information_schema.views
	WHERE table_schema = current_schema()
)`
}

func (d pgsqlDialect) SchemaTableName() string {
	return `
(
	SELECT
		table_name AS name,
		'table' AS type
	FROM information_schema.tables
	WHERE table_schema = current_schema() AND table_type = 'BASE TABLE'
	UNION ALL
	SELECT
		table_name AS name,
		'view' AS type
	FROM information_schema.views
	WHERE table_schema = current_schema()
)`
}

func (d pgsqlDialect) OptimizeSQL() string {
	return ""
}

func (d pgsqlDialect) WalCheckpointSQL() string {
	return ""
}

func (d pgsqlDialect) DateHourExpr(column string) string {
	return fmt.Sprintf("to_char(date_trunc('hour', %s::timestamp), 'YYYY-MM-DD HH24:00:00')", column)
}

func (d pgsqlDialect) JSONEach(column string) string {
	return fmt.Sprintf(
		"jsonb_array_elements_text(CASE WHEN [[%s]] IS NULL OR [[%s]] = '' THEN '[]'::jsonb WHEN LEFT(TRIM([[%s]]), 1) = '[' THEN [[%s]]::jsonb ELSE jsonb_build_array(to_jsonb([[%s]])) END)",
		column, column, column, column, column,
	)
}

func (d pgsqlDialect) JSONArrayLength(column string) string {
	return fmt.Sprintf(
		"jsonb_array_length(CASE WHEN [[%s]] IS NULL OR [[%s]] = '' THEN '[]'::jsonb WHEN LEFT(TRIM([[%s]]), 1) = '[' THEN [[%s]]::jsonb ELSE jsonb_build_array(to_jsonb([[%s]])) END)",
		column, column, column, column, column,
	)
}

func (d pgsqlDialect) JSONExtract(column string, path string) string {
	parts := []string{}
	if strings.TrimSpace(path) != "" {
		path = strings.TrimPrefix(path, ".")
		parts = strings.Split(path, ".")
	}

	pathArgs := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		pathArgs = append(pathArgs, fmt.Sprintf("'%s'", strings.ReplaceAll(p, "'", "''")))
	}

	args := strings.Join(pathArgs, ", ")
	if args != "" {
		args = ", " + args
	}

	return fmt.Sprintf(
		"(CASE WHEN [[%s]] IS NULL OR [[%s]] = '' THEN NULL WHEN LEFT(TRIM([[%s]]), 1) IN ('[', '{') THEN jsonb_extract_path([[%s]]::jsonb%s) ELSE jsonb_extract_path(jsonb_build_object('pb', to_jsonb([[%s]])), 'pb'%s) END)",
		column, column, column, column, args, column, args,
	)
}
