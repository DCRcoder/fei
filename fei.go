package fei

import (
	"context"
	"fmt"
)

// Config 配置
type Config struct {
	Driver       string
	MasterAddr   string
	SlavesAddr   []string
	MaxIdleConns int
	MaxOpenConns int
}

// Engine orm engine define
type Engine struct {
	*DB
	*Config
}

// NewEngine return engine
func NewEngine(driverName, dataSourceName string) (*Engine, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	e := &Engine{db, nil}
	return e, nil
}

// NewEngineWithMS return engine with master and slaves
func NewEngineWithMS(driverName, masterAddr string, slavesAddr []string) (*Engine, error) {
	db, err := OpenMasterAndSlaves(driverName, masterAddr, slavesAddr)
	if err != nil {
		return nil, err
	}
	e := &Engine{db, nil}
	return e, nil
}

// New return engine instance
func New(cfg *Config) (*Engine, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg empty")
	}
	var e *Engine
	if cfg.SlavesAddr != nil && len(cfg.SlavesAddr) != 0 {
		db, err := OpenMasterAndSlaves(cfg.Driver, cfg.MasterAddr, cfg.SlavesAddr)
		if err != nil {
			return nil, err
		}
		e = &Engine{db, cfg}
	} else {
		db, err := Open(cfg.Driver, cfg.MasterAddr)
		if err != nil {
			return nil, err
		}
		e = &Engine{db, cfg}
	}
	e.SetMaxIdleConns(cfg.MaxIdleConns)
	e.SetMaxOpenConns(cfg.MaxOpenConns)
	return e, nil
}

// NewSessionCtx return new sessiont instance with ctx
func (e *Engine) NewSessionCtx(ctx context.Context) *Session {
	return &Session{
		db:        e.DB,
		ctx:       ctx,
		statement: &Statement{},
	}
}

// NewSession return new session instance
func (e *Engine) NewSession() *Session {
	return e.NewSessionCtx(context.Background())
}
