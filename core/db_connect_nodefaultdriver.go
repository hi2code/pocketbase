//go:build no_default_driver

package core

import (
	"strings"

	"github.com/pocketbase/dbx"
)

func DefaultDBConnect(dbPath string) (*dbx.DB, error) {
	panic("DBConnect config option must be set when the no_default_driver tag is used!")
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
