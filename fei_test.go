package fei

import (
	"testing"

	"github.com/alecthomas/assert"
	_ "github.com/go-sql-driver/mysql"
)

type CodeBook struct {
	Name      string
	ID        int64
	Password  string
	Remarks   *string
	CreatedAt string
	UpdatedAt string
}

func (c *CodeBook) TableName() string {
	return "codebook"
}

func TestCount(t *testing.T) {
	engine, err := NewEngine("mysql", "root@/test?charset=utf8")
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	count, err := engine.NewSession().Select().From("codebook").Where(Eq{"name": "liubin"}).Count()
	assert.Equal(t, err, nil)
	assert.Equal(t, count, int64(3))
}

func TestFindOne(t *testing.T) {
	c := CodeBook{}
	engine, err := NewEngine("mysql", "root@/test?charset=utf8")
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	err = engine.NewSession().Select().Where(Eq{"name": "liubin"}).FindOne(&c)
	assert.Equal(t, err, nil)
	assert.Equal(t, c.Name, "liubin")
}

func TestFindAll(t *testing.T) {
	c := make([]*CodeBook, 0)
	engine, err := NewEngine("mysql", "root@/test?charset=utf8")
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	err = engine.NewSession().Select().Where(Eq{"name": "liubin"}).OrderBy("id desc").FindAll(&c)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(c), 2)
	assert.Equal(t, c[1].Name, "liubin")
	assert.Equal(t, c[1].Password, "qingning")
}
