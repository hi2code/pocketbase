package core_test

import "regexp"

var (
	testSelectTablePattern         = regexp.MustCompile("SELECT `([^`]+)`\\.\\*")
	testSelectDistinctTablePattern = regexp.MustCompile("SELECT DISTINCT `([^`]+)`\\.\\*")
	testFromTablePattern           = regexp.MustCompile("FROM `([^`]+)`")
	testLeftJoinTablePattern       = regexp.MustCompile("LEFT JOIN `([^`]+)`")
)

func normalizeCollectionTableRefsForTests(sql string) string {
	sql = testSelectDistinctTablePattern.ReplaceAllString(sql, "SELECT DISTINCT {{$1}}.*")
	sql = testSelectTablePattern.ReplaceAllString(sql, "SELECT {{$1}}.*")
	sql = testFromTablePattern.ReplaceAllString(sql, "FROM {{$1}}")
	sql = testLeftJoinTablePattern.ReplaceAllString(sql, "LEFT JOIN {{$1}}")
	return sql
}
