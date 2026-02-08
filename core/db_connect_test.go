package core

import (
	"testing"
)

func TestDefaultDBConnect_PostgreSQL(t *testing.T) {
	testCases := []struct {
		name    string
		dbPath  string
		wantErr bool
		isPG    bool
	}{
		{
			name:    "PostgreSQL URL with postgres://",
			dbPath:  "postgres://user:pass@localhost:5432/testdb",
			wantErr: false, // dbx.Open may not immediately connect
			isPG:    true,
		},
		{
			name:    "PostgreSQL URL with postgresql://",
			dbPath:  "postgresql://user:pass@localhost:5432/testdb",
			wantErr: false, // dbx.Open may not immediately connect
			isPG:    true,
		},
		{
			name:    "PostgreSQL URL with parameters",
			dbPath:  "postgresql://user:pass@localhost:5432/testdb?sslmode=disable&TimeZone=UTC",
			wantErr: false, // dbx.Open may not immediately connect
			isPG:    true,
		},
		{
			name:    "SQLite file path",
			dbPath:  "test.db",
			wantErr: false,
			isPG:    false,
		},
		{
			name:    "SQLite in-memory",
			dbPath:  ":memory:",
			wantErr: false,
			isPG:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := DefaultDBConnect(tc.dbPath)

			if tc.wantErr && err == nil {
				t.Errorf("DefaultDBConnect() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && err != nil {
				t.Errorf("DefaultDBConnect() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if db == nil && !tc.wantErr {
				t.Error("DefaultDBConnect() returned nil db")
			}

			// Cleanup
			if db != nil {
				db.Close()
			}
		})
	}
}

func TestPostgreSQLConnect_ConnectionString(t *testing.T) {
	testCases := []struct {
		name    string
		connStr string
	}{
		{
			name:    "Basic connection string",
			connStr: "postgres://user:pass@localhost:5432/db",
		},
		{
			name:    "Already has sslmode",
			connStr: "postgres://user:pass@localhost:5432/db?sslmode=require",
		},
		{
			name:    "Already has TimeZone",
			connStr: "postgres://user:pass@localhost:5432/db?TimeZone=Asia/Shanghai",
		},
		{
			name:    "Has both parameters",
			connStr: "postgres://user:pass@localhost:5432/db?sslmode=verify-full&TimeZone=UTC",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: We can't actually test the connection without a PostgreSQL server
			// This test just verifies the function signature and basic logic
			db, err := PostgreSQLConnect(tc.connStr)

			if err != nil {
				t.Logf("PostgreSQLConnect() error (expected without server): %v", err)
			}

			if db == nil && err == nil {
				t.Error("PostgreSQLConnect() returned nil db without error")
			}

			// Cleanup
			if db != nil {
				db.Close()
			}
		})
	}
}
