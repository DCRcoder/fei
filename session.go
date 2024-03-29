package fei

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
)

// Session db conn session
type Session struct {
	e                      error
	db                     *DB
	ctx                    context.Context
	statement              *Statement
	useMaster              bool
	logger                 Logger
	isAutoCommit           bool
	hasCommittedOrRollback bool
	tx                     *sql.Tx
	explainModel           bool
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
	if scanner.entityPointer.Kind() != reflect.Struct {
		return FindOneExpectStruct
	}
	defer scanner.Close()
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}

	if s.explainModel {
		s.explain()
	}

	s.initCtx()
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session FindOne] sql: %s, args: %v", sql, args)
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
	if scanner.entityPointer.Kind() != reflect.Slice {
		return FindAllExpectSlice
	}
	defer scanner.Close()
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}

	s.initCtx()
	if s.explainModel {
		s.explain()
	}
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session FindAll] sql: %s, args: %v", sql, args)
	rows, err := s.QueryContext(s.ctx, sql, args...)
	if err != nil {
		return err
	}
	scanner.SetRows(rows)
	return scanner.Convert()
}

// Insert create new record
func (s *Session) Insert(dest interface{}) (int64, error) {
	s.initStatemnt()
	s.statement.Insert()
	scanner, err := NewScanner(dest)
	if err != nil {
		return 0, err
	}
	defer scanner.Close()
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}
	insertFields := make([]string, 0)
	for n, f := range scanner.Model.Fields {
		if !f.IsReadOnly {
			insertFields = append(insertFields, n)
		}
	}
	s.Columns(insertFields...)
	if scanner.entityPointer.Kind() == reflect.Slice {
		for i := 0; i < scanner.entityPointer.Len(); i++ {
			val := make([]interface{}, 0)
			sub := scanner.entityPointer.Index(i)
			if sub.Kind() == reflect.Ptr {
				subElem := sub.Elem()
				for _, fn := range insertFields {
					f, ok := scanner.Model.Fields[fn]
					if !ok {
						continue
					}
					fv := subElem.Field(f.idx)
					val = append(val, fv.Interface())
				}

			} else {
				for _, fn := range insertFields {
					f, ok := scanner.Model.Fields[fn]
					if !ok {
						continue
					}
					fv := sub.Field(f.idx)
					val = append(val, fv.Interface())
				}
			}
			s.statement.Values(val)
		}

	} else if scanner.entityPointer.Kind() == reflect.Struct {
		val := make([]interface{}, 0)
		for _, fn := range insertFields {
			f, ok := scanner.Model.Fields[fn]
			if !ok {
				continue
			}
			fv := scanner.entityPointer.Field(f.idx)
			val = append(val, fv.Interface())
		}
		s.statement.Values(val)
	} else {
		return 0, InsertExpectSliceOrStruct
	}
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return 0, err
	}
	s.logger.Debugf("[Session Insert] sql: %s, args: %v", sql, args)
	s.initCtx()
	sResult, err := s.ExecContext(s.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return sResult.RowsAffected()
}

// Update update one record
func (s *Session) Update(dest interface{}) (int64, error) {
	s.initStatemnt()
	s.statement.Update()
	scanner, err := NewScanner(dest)
	if err != nil {
		return 0, err
	}
	defer scanner.Close()
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}
	updateFields := make([]string, 0)
	pks := make([]interface{}, 0)
	primaryKey := ""
	for n, f := range scanner.Model.Fields {
		if !f.IsReadOnly && !f.IsPrimaryKey {
			updateFields = append(updateFields, n)
		}
		if f.IsPrimaryKey {
			primaryKey = n
		}
	}
	if primaryKey == "" {
		return 0, ModelMustHavePrimaryKey
	}
	s.Columns(updateFields...)
	if scanner.entityPointer.Kind() == reflect.Slice {
		for i := 0; i < scanner.entityPointer.Len(); i++ {
			val := make([]interface{}, 0)
			sub := scanner.entityPointer.Index(i)
			if sub.Kind() == reflect.Ptr {
				subElem := sub.Elem()
				for _, fn := range updateFields {
					f, ok := scanner.Model.Fields[fn]
					if !ok {
						continue
					}
					fv := subElem.Field(f.idx)
					val = append(val, fv.Interface())
				}
				primaryF, _ := scanner.Model.Fields[primaryKey]
				fv := subElem.Field(primaryF.idx)
				pks = append(pks, fv.Interface())
			} else {
				for _, fn := range updateFields {
					f, ok := scanner.Model.Fields[fn]
					if !ok {
						continue
					}
					fv := sub.Field(f.idx)
					val = append(val, fv.Interface())
				}
				primaryF, _ := scanner.Model.Fields[primaryKey]
				fv := sub.Field(primaryF.idx)
				pks = append(pks, fv.Interface())
			}
			s.statement.Values(val)
		}

	} else if scanner.entityPointer.Kind() == reflect.Struct {
		val := make([]interface{}, 0)
		for _, fn := range updateFields {
			f, ok := scanner.Model.Fields[fn]
			if !ok {
				continue
			}
			fv := scanner.entityPointer.Field(f.idx)
			val = append(val, fv.Interface())
		}
		primaryF, _ := scanner.Model.Fields[primaryKey]
		fv := scanner.entityPointer.Field(primaryF.idx)
		pks = append(pks, fv.Interface())
		s.statement.Values(val)
	} else {
		return 0, UpdateExpectSliceOrStruct
	}
	s.Where(Eq{scanner.Model.PkName: pks})
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return 0, err
	}
	s.logger.Debugf("[Session Update] sql: %s, args: %v", sql, args)
	s.initCtx()
	sResult, err := s.ExecContext(s.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return sResult.RowsAffected()
}

// UpdateRow updateRow method without model
func (s *Session) UpdateRow(param map[string]interface{}) (int64, error) {
	if param == nil {
		return 0, UpdateRowParamMustHaveValue
	}
	s.statement.stType = UpdateStatement
	if len(s.statement.conditions) == 0 || s.statement.table == "" {
		return 0, UpdateRowMustWithConditionAndTableName
	}
	updateFields := make([]string, 0)
	val := make([]interface{}, 0)
	for k, v := range param {
		updateFields = append(updateFields, k)
		val = append(val, v)
	}
	s.Columns(updateFields...)
	s.statement.Values(val)
	sql, args, err := s.statement.ToSQL()
	s.logger.Debugf("[Session UpdateRow] sql: %s, args: %v", sql, args)
	s.initCtx()
	sResult, err := s.ExecContext(s.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return sResult.RowsAffected()
}

// Delete delete one record
func (s *Session) Delete(dest interface{}) (int64, error) {
	s.initStatemnt()
	s.statement.Delete()
	scanner, err := NewScanner(dest)
	if err != nil {
		return 0, err
	}
	defer scanner.Close()
	if s.statement.table == "" {
		s.statement.From(scanner.GetTableName())
	}
	pks := make([]interface{}, 0)
	if scanner.Model.PkName == "" {
		return 0, ModelMissingPrimaryKey
	}
	if scanner.entityPointer.Kind() == reflect.Slice {
		for i := 0; i < scanner.entityPointer.Len(); i++ {
			sub := scanner.entityPointer.Index(i)
			if sub.Kind() == reflect.Ptr {
				pks = append(pks, sub.Elem().Field(scanner.Model.PkIdx).Interface())
			} else {
				pks = append(pks, sub.Field(scanner.Model.PkIdx).Interface())
			}
		}
	} else if scanner.entityPointer.Kind() == reflect.Struct {
		pks = append(pks, scanner.entityPointer.Field(scanner.Model.PkIdx).Interface())
	} else {
		return 0, DeleteExpectSliceOrStruct
	}
	s.Where(Eq{scanner.Model.PkName: pks})
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		return 0, err
	}
	s.logger.Debugf("[Session Delete] sql: %s, args: %v", sql, args)
	s.initCtx()
	sResult, err := s.ExecContext(s.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return sResult.RowsAffected()
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
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) initStatemnt() {
	if s.statement == nil {
		s.statement = &Statement{}
	}
}

func (s *Session) EnableExplain(flag bool) *Session {
	s.explainModel = flag
	return s
}

// Select select columns default "*"
func (s *Session) Select(columns ...string) *Session {
	s.initStatemnt()
	s.statement.Select(columns...)
	return s
}

func (s *Session) UseIndexs(idx ...string) *Session {
	s.initStatemnt()
	s.statement.UseIndexs(idx...)
	return s
}

func (s *Session) ForceIndexs(idx ...string) *Session {
	s.initStatemnt()
	s.statement.ForceIndexs(idx...)
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
	if s.tx != nil {
		return s.tx.QueryRow(query, args...)
	}
	if s.useMaster {
		return s.db.Master().QueryRow(query, args...)
	}
	return s.db.Slave().QueryRow(query, args...)
}

// Query use Query with session config
func (s *Session) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.tx != nil {
		return s.tx.Query(query, args...)
	}
	if s.useMaster {
		return s.db.Master().Query(query, args...)
	}
	return s.db.Slave().Query(query, args...)
}

// QueryContext use QueryContext with session config
func (s *Session) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if s.tx != nil {
		return s.tx.QueryContext(ctx, query, args...)
	}
	if s.useMaster {
		return s.db.Master().QueryContext(ctx, query, args...)
	}
	return s.db.Slave().QueryContext(ctx, query, args...)
}

// QueryRawContext use QueryRawContext with session config
func (s *Session) QueryRawContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if s.tx != nil {
		return s.tx.QueryRowContext(ctx, query, args...)
	}
	if s.useMaster {
		return s.db.Master().QueryRowContext(ctx, query, args...)
	}
	return s.db.Slave().QueryRowContext(ctx, query, args...)
}

// Exec execute
func (s *Session) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.tx != nil {
		return s.tx.Exec(query, args...)
	}
	return s.db.Master().Exec(query, args...)
}

// ExecContext execute with context
func (s *Session) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if s.tx != nil {
		return s.tx.ExecContext(ctx, query, args...)
	}
	return s.db.Master().ExecContext(ctx, query, args...)
}

// Begin begin transaction
func (s *Session) Begin() error {
	tx, err := s.db.Master().Begin()
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session Begin] begin transaction")
	// 重置 transaction 状态
	s.hasCommittedOrRollback = false
	s.isAutoCommit = false
	s.tx = tx
	return nil
}

// BeginTx begin transaction with opts
func (s *Session) BeginTx(opts *sql.TxOptions) error {
	tx, err := s.db.Master().BeginTx(s.ctx, opts)
	if err != nil {
		return err
	}
	s.logger.Debugf("[Session Begin] begin transaction opts: %v", opts)
	// 重置 transaction 状态
	s.hasCommittedOrRollback = false
	s.isAutoCommit = false
	s.tx = tx
	return nil
}

// RollBack  transaction rollback
func (s *Session) RollBack() error {
	if !s.isAutoCommit && !s.hasCommittedOrRollback {
		err := s.tx.Rollback()
		if err != nil {
			return err
		}
		s.hasCommittedOrRollback = true
	}
	return nil
}

// Commit transaction commit
func (s *Session) Commit() error {
	if !s.isAutoCommit && !s.hasCommittedOrRollback {
		err := s.tx.Commit()
		if err != nil {
			return err
		}
		s.hasCommittedOrRollback = true
	}
	return nil
}

// Transaction  transaction with autoCommit
func (s *Session) Transaction(f func(*Session) (interface{}, error)) (interface{}, error) {
	err := s.Begin()
	if err != nil {
		return nil, err
	}
	d, err := f(s)
	if err != nil {
		s.RollBack()
	} else {
		s.Commit()
	}
	return d, err
}

// TransactionTx  transactionTx with autoCommit
func (s *Session) TransactionTx(f func(*Session) (interface{}, error), opts *sql.TxOptions) (interface{}, error) {
	err := s.BeginTx(opts)
	if err != nil {
		return nil, err
	}
	d, err := f(s)
	if err != nil {
		s.RollBack()
	} else {
		s.Commit()
	}
	return d, err
}

func (s *Session) explain() {
	s.statement.EnableExplain(true)
	defer s.statement.EnableExplain(false)
	sql, args, err := s.statement.ToSQL()
	if err != nil {
		s.logger.Debugf("[Session FindAll] explian tosql error sql: %s, args: %v error: %v", sql, args, err)
		return
	}
	rows, err := s.QueryContext(s.ctx, sql, args...)
	if err != nil {
		s.logger.Debugf("[Session FindAll] explian query error sql: %s, args: %v error: %v", sql, args, err)
		return
	}
	for rows.Next() {
		explain := &ExplainModel{}
		err := rows.Scan(&explain.ID, &explain.SelectType, &explain.Table, &explain.Partitions, &explain.Type, &explain.PossibleKeys, &explain.Key, &explain.KeyLen, &explain.Ref, &explain.Rows, &explain.Filtered, &explain.Extra)
		if err != nil {
			s.logger.Debugf("[Session FindAll] explian model scan error sql: %s, args: %v error: %v", sql, args, err)
			return
		}
		b, _ := json.Marshal(explain)
		s.logger.Debugf("[Session FindAll] explian %s", string(b))
	}
}
