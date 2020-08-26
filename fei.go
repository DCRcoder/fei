package fei

import (
	"context"
)

// Config 配置
type Config struct {
	Driver       string
	MasterAddr   string
	SlavesAddr   []string
	MaxIdleConns int
	MaxOpenConns int
	Logger       Logger
	LogLevel     LogLevel
}

// Engine orm engine define
type Engine struct {
	*DB
	*Config
	Logger Logger
}

// NewEngine return engine
func NewEngine(driverName, dataSourceName string) (*Engine, error) {
	e := &Engine{nil, nil, NewFlogger()}
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	e.DB = db
	return e, nil
}

// NewEngineWithMS return engine with master and slaves
func NewEngineWithMS(driverName, masterAddr string, slavesAddr []string) (*Engine, error) {
	e := &Engine{nil, nil, NewFlogger()}
	db, err := OpenMasterAndSlaves(driverName, masterAddr, slavesAddr)
	if err != nil {
		return nil, err
	}
	e.DB = db
	return e, nil
}

// New return engine instance
func New(cfg *Config) (*Engine, error) {
	if cfg == nil {
		return nil, CFBNotAllowEmpty
	}
	e := &Engine{nil, cfg, nil}
	if cfg.Logger != nil {
		e.Logger = cfg.Logger
	} else if cfg.LogLevel != LogUnknown {
		e.Logger = NewFlogger(cfg.LogLevel)
	} else {
		e.Logger = NewFlogger()
	}
	if cfg.SlavesAddr != nil && len(cfg.SlavesAddr) != 0 {
		e.Logger.Debugf("[New Engine] driver: %s, masterAddr: %s, slaveAddr: %v", cfg.Driver, cfg.MasterAddr, cfg.SlavesAddr)
		db, err := OpenMasterAndSlaves(cfg.Driver, cfg.MasterAddr, cfg.SlavesAddr)
		if err != nil {
			return nil, err
		}
		e.DB = db
	} else {
		e.Logger.Debugf("[New Engine] driver: %s, masterAddr: %s", cfg.Driver, cfg.MasterAddr)
		db, err := Open(cfg.Driver, cfg.MasterAddr)
		if err != nil {
			return nil, err
		}
		e.DB = db
	}
	e.SetLogLevel(LogError)
	e.SetMaxIdleConns(cfg.MaxIdleConns)
	e.SetMaxOpenConns(cfg.MaxOpenConns)
	return e, nil
}

// SetLogger set loggger
func (e *Engine) SetLogger(logger Logger) {
	if logger != nil {
		e.Logger = logger
	}
}

// SetLogLevel set logger level
func (e *Engine) SetLogLevel(level LogLevel) {
	e.Logger.SetLogLevel(level)
}

// NewSessionCtx return new sessiont instance with ctx
func (e *Engine) NewSessionCtx(ctx context.Context) *Session {
	return &Session{
		db:                     e.DB,
		ctx:                    ctx,
		statement:              &Statement{},
		logger:                 e.Logger,
		isAutoCommit:           true,
		hasCommittedOrRollback: false,
		tx:                     nil,
	}
}

// NewSession return new session instance
func (e *Engine) NewSession() *Session {
	return e.NewSessionCtx(nil)
}

// Close Engine close
func (e *Engine) Close() {
	e.DB.Close()
}
