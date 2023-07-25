package checks

import (
	"database/sql"
)

func NewSQLCheck(sql *sql.DB) Check {
	return &sqlCheck{
		sql: sql,
	}
}

type sqlCheck struct {
	sql *sql.DB
}

func (c *sqlCheck) Pass() bool {
	if c.sql == nil {
		return false
	}
	return c.sql.Ping() == nil
}

func (c *sqlCheck) Name() string {
	return "Mysql"
}
