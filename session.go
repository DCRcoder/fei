package fei

import (
	"context"
	"database/sql"
)

// Session db conn session
type Session struct {
	e         error
	db        *DB
	ctx       context.Context
	statement *Statement
	useMaster bool
	logger    Logger
}

// UseMaster enable use master
func (s *Session) UseMaster() *Session {
	s.useMaster = true
	return s
}

// FindOne get one result
func (s *Session) FindOne(dest interface{}) error {
	return nil
}

// FindAll get all result
func (s *Session) FindAll(dest interface{}) error {
	return nil
}

// Count return query count
func (s *Session) Count() (int64, error) {
	s.initStatemnt()
	s.Columns("count(*)")
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return 0, err
	}
	s.logger.Debugf("[Session Count] sql: %s, args: %v", sql, args)
	var count int64
	row := s.QueryRow(sql, args...)
	row.Scan(&count)
	return count, nil
}

func (s *Session) initStatemnt() {
	if s.statement == nil {
		s.statement = &Statement{}
	}
}

// Select select columns default "*"
func (s *Session) Select(columns ...string) *Session {
	s.initStatemnt()
	s.statement.Select(columns...)
	return s
}

// Columns set sql columns atttentio Columns will reset st.columns
func (s *Session) Columns(columns ...string) *Session {
	s.initStatemnt()
	s.statement.Columns(columns...)
	return s
}

// From set select table
func (s *Session) From(table string) *Session {
	s.initStatemnt()
	s.statement.From(table)
	return s
}

// Where set conditions
func (s *Session) Where(expr ...interface{}) *Session {
	s.initStatemnt()
	s.statement.Where(expr...)
	return s
}

// QueryRow use QueryRow with session config
func (s *Session) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.useMaster {
		return s.db.Master().QueryRowContext(s.ctx, query, args...)
	}
	return s.db.Slave().QueryRowContext(s.ctx, query, args...)
}

// Query use Query with session config
func (s *Session) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.useMaster {
		return s.db.Master().QueryContext(s.ctx, query, args...)
	}
	return s.db.Slave().QueryContext(s.ctx, query, args...)
}
