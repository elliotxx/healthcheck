package checks

import (
	"database/sql"
)

// NewSQLCheck creates a new SQL health check.
func NewSQLCheck(sql *sql.DB) Check {
	return &sqlCheck{
		sql: sql,
	}
}

// sqlCheck represents an SQL health check.
type sqlCheck struct {
	sql *sql.DB
}

// Pass checks if the SQL database is reachable.
func (c *sqlCheck) Pass() bool {
	if c.sql == nil {
		return false
	}
	return c.sql.Ping() == nil
}

// Name returns the name of the SQL health check.
func (c *sqlCheck) Name() string {
	return "Mysql"
}
