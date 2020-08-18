package fei

import "fmt"

var (
	CFBNotAllowEmpty     = fmt.Errorf("config not allow empty")
	StatementTableNotSet = fmt.Errorf("statement table not set")
	StatementTypeNotSet  = fmt.Errorf("statement type not set")
)
