package dbutils

// DmDialect follows mysql-compatible SQL behavior.
//
// DM is configured in this project to run in MySQL compatibility mode, so we
// intentionally reuse mysqlDialect expressions for cross-db consistency.
type DmDialect struct{}

func init() {
	RegisterDMDialect("dm")
}

func RegisterDMDialect(driverNames ...string) {
	if len(driverNames) == 0 {
		driverNames = []string{"dm"}
	}

	for _, name := range driverNames {
		RegisterDialect(name, DmDialect{})
	}
}

func (d DmDialect) Name() string {
	return "dm"
}

func (d DmDialect) TableColumnsSQL() string {
	return mysqlDialect{}.TableColumnsSQL()
}

func (d DmDialect) TableInfoSQL() string {
	return mysqlDialect{}.TableInfoSQL()
}

func (d DmDialect) MasterTableName() string {
	return mysqlDialect{}.MasterTableName()
}

func (d DmDialect) SchemaTableName() string {
	return mysqlDialect{}.SchemaTableName()
}

func (d DmDialect) OptimizeSQL() string {
	return mysqlDialect{}.OptimizeSQL()
}

func (d DmDialect) WalCheckpointSQL() string {
	return mysqlDialect{}.WalCheckpointSQL()
}

func (d DmDialect) DateHourExpr(column string) string {
	return mysqlDialect{}.DateHourExpr(column)
}

func (d DmDialect) JSONEach(column string) string {
	return mysqlDialect{}.JSONEach(column)
}

func (d DmDialect) JSONArrayLength(column string) string {
	return mysqlDialect{}.JSONArrayLength(column)
}

func (d DmDialect) JSONExtract(column, path string) string {
	return mysqlDialect{}.JSONExtract(column, path)
}
