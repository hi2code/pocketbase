package dbutils

import "fmt"

// JSONEach returns JSON_EACH SQLite string expression with
// some normalizations for non-json columns.
func JSONEach(column string) string {
	return GetDialect().JSONEach(column)
}

// JSONEachParam returns a table expression that expands a JSON array parameter
// into rows with a single "value" column.
func JSONEachParam(placeholder string) string {
	switch GetDialect().Name() {
	case "pg":
		return fmt.Sprintf("jsonb_array_elements_text((:%s)::jsonb)", placeholder)
	case "mysql", "dm":
		return fmt.Sprintf(
			"(SELECT JSON_UNQUOTE(JSON_EXTRACT({:%s}, CONCAT('$[', _pb_seq.n, ']'))) AS value FROM _pb_seq WHERE _pb_seq.n < JSON_LENGTH({:%s}))",
			placeholder,
			placeholder,
		)
	default:
		return fmt.Sprintf("json_each({:%s})", placeholder)
	}
}

// JSONArrayLength returns JSON_ARRAY_LENGTH SQLite string expression
// with some normalizations for non-json columns.
//
// It works with both json and non-json column values.
//
// Returns 0 for empty string or NULL column values.
func JSONArrayLength(column string) string {
	return GetDialect().JSONArrayLength(column)
}

// JSONExtract returns a JSON_EXTRACT SQLite string expression with
// some normalizations for non-json columns.
func JSONExtract(column string, path string) string {
	return GetDialect().JSONExtract(column, path)
}
