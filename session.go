package fei

import (
	"context"
)

// Session db conn session
type Session struct {
	e         error
	db        *DB
	ctx       context.Context
	statement *Statement
	useMaster bool
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
	var count int64
	row := s.db.Slave().QueryRow(sql, args...)
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
