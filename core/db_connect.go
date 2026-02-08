//go:build !no_default_driver

package core

import (
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func DefaultDBConnect(dbPath string) (*dbx.DB, error) {
	// Check if this is a PostgreSQL connection string
	if strings.HasPrefix(dbPath, "postgres://") || strings.HasPrefix(dbPath, "postgresql://") {
		return PostgreSQLConnect(dbPath)
	}

	// Default to SQLite
	// Note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection.
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"

	db, err := dbx.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// PostgreSQLConnect creates a PostgreSQL database connection using pgx driver
func PostgreSQLConnect(connStr string) (*dbx.DB, error) {
	// Parse connection string to ensure it's valid
	// Add default parameters if not present
	if !strings.Contains(connStr, "sslmode=") {
		if strings.Contains(connStr, "?") {
			connStr += "&sslmode=disable"
		} else {
			connStr += "?sslmode=disable"
		}
	}

	// Add timezone if not present
	if !strings.Contains(connStr, "TimeZone=") {
		if strings.Contains(connStr, "?") {
			connStr += "&TimeZone=UTC"
		} else {
			connStr += "?TimeZone=UTC"
		}
	}

	// pgx uses "pgx" as driver name
	db, err := dbx.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
