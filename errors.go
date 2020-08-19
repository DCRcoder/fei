package fei

import (
	"errors"
)

var (
	CFBNotAllowEmpty            = errors.New("config not allow empty")
	StatementTableNotSet        = errors.New("statement table not set")
	StatementTypeNotSet         = errors.New("statement type not set")
	ScannerRowsPointerNil       = errors.New("Scanner rows could not be nil pointer")
	ScannerEntityNeedCanSet     = errors.New("Entity need can set")
	ScannerEntiryTypeNotSupport = errors.New("Scanner Entity not support. it should be pointer or slice")
)
