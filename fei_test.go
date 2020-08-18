package fei

import (
	"testing"

	"github.com/alecthomas/assert"
	_ "github.com/go-sql-driver/mysql"
)

func TestFindOne(t *testing.T) {
	engine, err := NewEngine("mysql", "root@/test?charset=utf8")
	assert.Equal(t, err, nil)
	count, err := engine.NewSession().Select().From("codebook").Where(Eq{"name": "laojun"}).Count()
	assert.Equal(t, err, nil)
	assert.Equal(t, count, int64(0))
}
