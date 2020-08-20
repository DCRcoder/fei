package fei

import (
	"database/sql"
	"fmt"
	"reflect"
)

const (
	feiColumnName = "fc_name"
)

// Scanner convert rows to entity
// Don't scan into interface{} but the type you would expect, the database/sql package converts the returned type for you then.
type Scanner struct {
	rows          *sql.Rows
	fields        []string
	entity        interface{}
	entityValue   reflect.Value
	entityPointer reflect.Value
	model         *Model
}

// Model describe table struct
type Model struct {
	TableName string
	Value     reflect.Value
	Fields    map[string]*Field
}

// Field describe table field
type Field struct {
	Name   string
	idx    int
	Column reflect.StructField
}

// NewModel return new model instanc
func NewModel(value reflect.Value) *Model {
	m := &Model{Value: value}
	_, ok := value.Type().MethodByName("TableName")
	if ok {
		vals := value.MethodByName("TableName").Call([]reflect.Value{})
		if len(vals) > 0 {
			switch vals[0].Kind() {
			case reflect.String:
				m.TableName = vals[0].String()
			}
		}
	}
	m.Fields = make(map[string]*Field)
	elem := value.Elem()
	for i := 0; i < elem.NumField(); i++ {
		df := elem.Type().Field(i)
		fieldName := ToSnakeCase(df.Name)
		if df.Tag.Get(feiColumnName) != "" {
			fieldName = df.Tag.Get(feiColumnName)
		}
		m.Fields[fieldName] = &Field{
			Name: fieldName,
			idx:  i,
		}
	}
	return m
}

// NewScanner return new scanner instance
func NewScanner(dest interface{}) (*Scanner, error) {
	entityValue := reflect.ValueOf(dest)
	fmt.Println(entityValue.Kind())
	s := &Scanner{
		entity:        dest,
		entityValue:   entityValue,
		entityPointer: reflect.Indirect(entityValue),
	}
	fmt.Println(s.entityPointer.Kind())
	if !s.entityPointer.CanSet() {
		return nil, ScannerEntityNeedCanSet
	}
	switch s.entityPointer.Kind() {
	case reflect.Slice:
		t := reflect.New(s.entityPointer.Type().Elem().Elem())
		s.model = NewModel(t)
	case reflect.Struct:
		s.model = NewModel(s.entityValue)
	default:
		return nil, ScannerEntiryTypeNotSupport
	}
	return s, nil
}

// Close close
func (sc *Scanner) Close() {
	if sc.rows != nil {
		sc.rows.Close()
	}
}

// SetRows set row
func (sc *Scanner) SetRows(rows *sql.Rows) {
	sc.rows = rows
}

// GetTableName try get table from dest
func (sc *Scanner) GetTableName() string {
	if sc.model != nil {
		return sc.model.TableName
	}
	return ""
}

// Convert scan rows to dest
func (sc *Scanner) Convert() error {
	srcValue := make([]interface{}, len(sc.fields))
	for i := 0; i < len(sc.fields); i++ {
		var v interface{}
		srcValue[i] = &v
	}
	if sc.rows == nil {
		return ScannerRowsPointerNil
	}
	fields, err := sc.rows.Columns()
	fmt.Println(fields)
	if err != nil {
		return err
	}
	sc.fields = fields
	switch sc.entityPointer.Kind() {
	case reflect.Slice:
		return sc.convertAll()
	case reflect.Struct:
		return sc.convertOne()
	default:
		return ScannerEntiryTypeNotSupport
	}
}

func (sc *Scanner) convertAll() error {
	dest := reflect.MakeSlice(sc.entityPointer.Type(), 0, 0)
	for sc.rows.Next() {
		srcValue := make([]interface{}, len(sc.fields))
		for i := 0; i < len(sc.fields); i++ {
			var v interface{}
			srcValue[i] = &v
		}
		err := sc.rows.Scan(srcValue...)
		if err != nil {
			return err
		}
		t := reflect.New(sc.entityPointer.Type().Elem().Elem())
		sc.SetEntity(srcValue, t.Elem())
		dest = reflect.Append(dest, t)
	}
	sc.entityPointer.Set(dest)
	return nil
}

func (sc *Scanner) convertOne() error {
	srcValue := make([]interface{}, len(sc.fields))
	for i := 0; i < len(sc.fields); i++ {
		var v interface{}
		srcValue[i] = &v
	}
	if sc.rows.Next() {
		err := sc.rows.Scan(srcValue...)
		if err != nil {
			return err
		}
		sc.SetEntity(srcValue, sc.entityPointer)
	}
	return nil
}

// SetEntity set entity
func (sc *Scanner) SetEntity(srcValue []interface{}, dest reflect.Value) error {
	tmpMap := make(map[string]interface{})
	for i := 0; i < len(sc.fields); i++ {
		field := sc.fields[i]
		value := srcValue[i]
		tmpMap[field] = value
	}
	for name, field := range sc.model.Fields {
		val, ok := tmpMap[name]
		if !ok {
			continue
		}
		ff := dest.Field(field.idx)
		rawVal := reflect.Indirect(reflect.ValueOf(val))
		if rawVal.Interface() == nil {
			continue
		}
		rawValueType := reflect.TypeOf(rawVal.Interface())
		vv := reflect.ValueOf(rawVal.Interface())
		fmt.Println(val, name, rawVal, rawValueType.Kind(), rawValueType, ff.Kind())
		switch ff.Kind() {
		case reflect.String:
			switch data := rawVal.Interface().(type) {
			case string:
				ff.SetString(data)
			case []uint8:
				ff.SetString(string(data))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch rawValueType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ff.SetInt(vv.Int())
			}
		case reflect.Float32, reflect.Float64:
			switch rawValueType.Kind() {
			case reflect.Float32, reflect.Float64:
				ff.SetFloat(vv.Float())
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			switch rawValueType.Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				ff.SetUint(vv.Uint())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ff.SetUint(uint64(vv.Int()))
			}
		}
	}
	return nil
}
