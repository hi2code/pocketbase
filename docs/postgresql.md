# PostgreSQL Support for PocketBase

PocketBase now includes experimental PostgreSQL support using the [pgx](https://github.com/jackc/pgx) driver. This allows you to use PostgreSQL as your database backend instead of the default SQLite.

## Quick Start

To use PostgreSQL with PocketBase, set your data directory to a PostgreSQL connection string:

```bash
# Using environment variable
PB_DATA_DIR="postgres://username:password@localhost:5432/pocketbase" ./pocketbase serve

# Using command line flag
./pocketbase serve --dir="postgres://username:password@localhost:5432/pocketbase"
```

## Connection String Format

PocketBase supports standard PostgreSQL connection strings:

```
postgres://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
```

### Common Examples:

1. **Local PostgreSQL with default port:**
   ```
   postgres://postgres:password@localhost:5432/pocketbase
   ```

2. **With SSL disabled:**
   ```
   postgres://user:pass@localhost:5432/dbname?sslmode=disable
   ```

3. **With specific timezone:**
   ```
   postgres://user:pass@localhost:5432/dbname?TimeZone=UTC
   ```

4. **With connection pool settings:**
   ```
   postgres://user:pass@localhost:5432/dbname?sslmode=disable&pool_max_conns=10
   ```

## Default Parameters

When not specified, PocketBase automatically adds:
- `sslmode=disable` (for local development)
- `TimeZone=UTC`

## Driver Information

PocketBase uses the [pgx](https://github.com/jackc/pgx) driver (v5) for PostgreSQL connectivity. pgx is a high-performance PostgreSQL driver that provides:
- Better performance than lib/pq
- Native support for PostgreSQL types
- Connection pooling
- Prepared statement caching

## Programmatic Usage

You can also configure PostgreSQL programmatically:

```go
package main

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	config := pocketbase.Config{
		DefaultDataDir: "postgres://username:password@localhost:5432/pocketbase",
	}
	
	app := pocketbase.NewWithConfig(config)
	
	if err := app.Start(); err != nil {
		panic(err)
	}
}
```

## Current Limitations

⚠️ **Important**: PostgreSQL support is experimental and has some limitations:

1. **SQLite-specific features**: Some PocketBase features rely on SQLite-specific SQL:
   - `PRAGMA` statements (used for optimization)
   - `sqlite_master`/`sqlite_schema` table queries
   - `_rowid_` column references
   - `PRAGMA_TABLE_INFO()` function

2. **Migration compatibility**: Existing SQLite migrations may need adaptation for PostgreSQL.

3. **Performance characteristics**: PostgreSQL has different performance characteristics compared to SQLite.

## Testing PostgreSQL Connection

To test if PostgreSQL connection works:

1. Start a PostgreSQL server (Docker example):
   ```bash
   docker run -d --name postgres-test \
     -e POSTGRES_PASSWORD=password \
     -e POSTGRES_DB=pocketbase \
     -p 5432:5432 \
     postgres:latest
   ```

2. Run PocketBase with PostgreSQL:
   ```bash
   PB_DATA_DIR="postgres://postgres:password@localhost:5432/pocketbase" ./pocketbase serve
   ```

## Troubleshooting

### Common Issues:

1. **Connection refused**: Ensure PostgreSQL is running and accessible
2. **Authentication failed**: Verify username/password
3. **Database doesn't exist**: Create the database first:
   ```sql
   CREATE DATABASE pocketbase;
   ```

### Error Messages:

- `pq: SSL is not enabled on the server`: Add `?sslmode=disable` to connection string
- `pq: database "pocketbase" does not exist`: Create the database first
- `dial tcp [::1]:5432: connect: connection refused`: PostgreSQL not running or wrong host/port

## Future Improvements

Planned enhancements for PostgreSQL support:

1. Database-agnostic query building
2. PostgreSQL-specific optimizations
3. Full migration compatibility
4. Connection pooling configuration
5. Read replicas support

## Contributing

If you encounter issues or want to improve PostgreSQL support, please:
1. Check existing issues on GitHub
2. Create detailed bug reports including:
   - PostgreSQL version
   - Connection string (with sensitive data redacted)
   - Error messages
   - Steps to reproduce

## License

PostgreSQL support is part of PocketBase and follows the same MIT license.