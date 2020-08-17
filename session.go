package fei

import "context"

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
