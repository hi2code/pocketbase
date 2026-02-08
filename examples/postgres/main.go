// PostgreSQL example for PocketBase using pgx driver
//
// To use PostgreSQL with PocketBase:
//
//  1. Set the data directory to a PostgreSQL connection string instead of a file path
//     Example connection strings:
//     - "postgres://username:password@localhost:5432/database_name"
//     - "postgresql://user:pass@localhost:5432/dbname?sslmode=disable&TimeZone=UTC"
//
//  2. You can set it via environment variable:
//     PB_DATA_DIR="postgres://username:password@localhost:5432/pocketbase" ./pocketbase serve
//
//  3. Or via command line flag:
//     ./pocketbase serve --dir="postgres://username:password@localhost:5432/pocketbase"
//
// Note: Full PostgreSQL compatibility requires additional work as PocketBase
//
//	currently has some SQLite-specific queries and features.
//	Uses pgx driver (github.com/jackc/pgx/v5) for PostgreSQL connectivity.
package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	// The PostgreSQL connection will be automatically detected
	// when the data directory is set to a PostgreSQL connection string

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
