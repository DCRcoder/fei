package fei

import (
	"database/sql"
	"fmt"
	"reflect"
)

const (
	feiColumnName = "fei_column"
)

// Scanner convert rows to model
type Scanner struct {
	rows          *sql.Rows
	fields        []string
	fieldTypes    []*sql.ColumnType
	entity        interface{}
	entityValue   reflect.Value
	entityPointer reflect.Value
}

// NewScanner return new instance
func NewScanner(dest interface{}) (*Scanner, error) {
	entityVal := reflect.ValueOf(dest)
	s := &Scanner{
		entity:        dest,
		entityValue:   entityVal,
		entityPointer: reflect.Indirect(entityVal),
	}
	if !s.entityPointer.CanSet() {
		return nil, ScannerEntityNeedCanSet
	}
	return s, nil
}

// SetRows set row
func (sc *Scanner) SetRows(rows *sql.Rows) {
	sc.rows = rows
}

// GetTableName try get table from dest
func (sc *Scanner) GetTableName() string {
	_, ok := sc.entityValue.Type().MethodByName("TableName")
	if ok {
		vals := sc.entityValue.MethodByName("TableName").Call([]reflect.Value{})
		if len(vals) > 0 {
			switch vals[0].Kind() {
			case reflect.String:
				return vals[0].String()
			}
		}
	}
	return ""
}

// Convert scan rows to des
func (sc *Scanner) Convert() error {
	if sc.rows == nil {
		return ScannerRowsPointerNil
	}
	fields, err := sc.rows.Columns()
	fmt.Println(fields)
	if err != nil {
		return err
	}
	sc.fields = fields
	fieldTypes, err := sc.rows.ColumnTypes()
	if err != nil {
		return err
	}
	sc.fieldTypes = fieldTypes
	fmt.Println(fieldTypes)
	switch sc.entityValue.Kind() {
	case reflect.Slice:
		fmt.Println("i m slice")
		return nil
	case reflect.Ptr:
		return sc.convertOne()
	default:
		return ScannerEntiryTypeNotSupport
	}
}

func (sc *Scanner) convertOne() error {
	srcValue := make([]interface{}, len(sc.fields))
	for i := 0; i < len(sc.fields); i++ {
		var v interface{}
		srcValue[i] = &v
	}
	if sc.rows.Next() {
		sc.rows.Scan(srcValue...)
		sc.SetEntity(srcValue)
		fmt.Println(sc.entity)
	}
	return nil
}

// SetEntity set entity
func (sc *Scanner) SetEntity(srcValue []interface{}) error {
	tmpMap := make(map[string]interface{})
	for i := 0; i < len(sc.fields); i++ {
		field := sc.fields[i]
		value := srcValue[i]
		tmpMap[field] = value
	}
	elem := sc.entityValue.Elem()
	for f := 0; f < elem.Type().NumField(); f++ {
		df := elem.Type().Field(f)
		fieldName := ToSnakeCase(df.Name)
		if df.Tag.Get(feiColumnName) != "" {
			fieldName = df.Tag.Get(feiColumnName)
		}
		val, ok := tmpMap[fieldName]
		if !ok {
			continue
		}
		ff := sc.entityPointer.Field(f)
		rawVal := reflect.Indirect(reflect.ValueOf(val))
		if rawVal.Interface() == nil {
			continue
		}
		rawValueType := reflect.TypeOf(rawVal.Interface())
		vv := reflect.ValueOf(rawVal.Interface())
		fmt.Println(val, fieldName, rawVal, rawValueType.Kind(), ff.Kind())
		switch ff.Kind() {
		case reflect.String:
			fmt.Println(vv.String())
			fmt.Println("i m string")
		case reflect.Int64:
			if rawValueType.Kind() == reflect.Int64 {
				var n = int64(vv.Int())
				ff.Set(reflect.ValueOf(n))
			}
		}
	}
	return nil
}
