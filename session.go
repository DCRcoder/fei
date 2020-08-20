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
	s.initStatemnt()
	s.Limit(1)
	scanner, err := NewScanner(dest)
	if err != nil {
		return err
	}
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session FindOne] sql: %s, args: %v", sql, args)
	s.initCtx()
	rows, err := s.QueryContext(s.ctx, sql, args...)
	if err != nil {
		return err
	}
	scanner.SetRows(rows)
	return scanner.Convert()
}

// FindAll get all result
func (s *Session) FindAll(dest interface{}) error {
	s.initStatemnt()
	scanner, err := NewScanner(dest)
	if err != nil {
		return err
	}
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session FindALL] sql: %s, args: %v", sql, args)
	s.initCtx()
	rows, err := s.QueryContext(s.ctx, sql, args...)
	if err != nil {
		return err
	}
	scanner.SetRows(rows)
	return scanner.Convert()
}

// Insert create new record
func (s *Session) Insert(model interface{}) (int64, error) {
	return 0, nil
}

// Update update one record
func (s *Session) Update(model interface{}) (int64, error) {
	return 0, nil
}

// Delete delete one record
func (s *Session) Delete(model interface{}) (int64, error) {
	return 0, nil
}

func (s *Session) initCtx() {
	if s.ctx == nil {
		s.ctx = context.Background()
	}
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
	s.initCtx()
	rows, err := s.QueryContext(s.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	rows.Scan(&count)
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

// Limit set limit
func (s *Session) Limit(limit uint64) *Session {
	s.statement.Limit(limit)
	return s
}

// Offset set limit
func (s *Session) Offset(offset uint64) *Session {
	s.statement.Offset(offset)
	return s
}

// OrderBy set order by
func (s *Session) OrderBy(orderby string) *Session {
	s.statement.OrderBy(orderby)
	return s
}

// QueryRow use QueryRow with session config
func (s *Session) QueryRow(query string, args ...interface{}) *sql.Row {
	if s.useMaster {
		return s.db.Master().QueryRow(query, args...)
	}
	return s.db.Slave().QueryRow(query, args...)
}

// Query use Query with session config
func (s *Session) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.useMaster {
		return s.db.Master().Query(query, args...)
	}
	return s.db.Slave().Query(query, args...)
}

// QueryContext use QueryContext with session config
func (s *Session) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if s.useMaster {
		return s.db.Master().QueryContext(ctx, query, args...)
	}
	return s.db.Slave().QueryContext(ctx, query, args...)
}

// QueryRawContext use QueryRawContext with session config
func (s *Session) QueryRawContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if s.useMaster {
		return s.db.Master().QueryRowContext(ctx, query, args...)
	}
	return s.db.Slave().QueryRowContext(ctx, query, args...)
}
