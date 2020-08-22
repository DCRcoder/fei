package fei

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
)

const dbAddr = "root@/test?charset=utf8"

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

type CodeBook struct {
	Name      string
	ID        int64 `fei:"pk,columnName=id"`
	Password  string
	Remarks   *string
	CreatedAt string    `fei:"readOnly"`
	UpdatedAt time.Time `fei:"readOnly"`
}

func (c *CodeBook) TableName() string {
	return "codebook"
}

func TestMain(m *testing.M) {
	var err error

	db, err := sql.Open("mysql", dbAddr)
	if err != nil {
		panic(err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("testdata"),
	)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

func TestCount(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	count, err := engine.NewSession().Select().From("codebook").Where(Eq{"name": "liubin"}).Count()
	assert.Equal(t, err, nil)
	assert.Equal(t, count, int64(3))
}

func TestFindOne(t *testing.T) {
	prepareTestDatabase()
	c := CodeBook{}
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	err = engine.NewSession().Select().Where(Eq{"name": "laojun"}).FindOne(&c)
	assert.Equal(t, err, nil)
	assert.Equal(t, c.Name, "laojun")
	assert.Equal(t, *c.Remarks, "qingning")
}

func TestFindOneColumn(t *testing.T) {
	prepareTestDatabase()
	c := CodeBook{}
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	err = engine.NewSession().Select("name").Where(Eq{"name": "laojun"}).FindOne(&c)
	assert.Equal(t, err, nil)
	assert.Equal(t, c.Name, "laojun")
	assert.Equal(t, c.ID, int64(0))
	assert.Equal(t, c.Remarks, (*string)(nil))
}

func TestFindAll(t *testing.T) {
	prepareTestDatabase()
	c := make([]*CodeBook, 0)
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	err = engine.NewSession().Select().Where(Eq{"name": "liubin"}).OrderBy("id desc").FindAll(&c)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(c), 3)
	assert.Equal(t, c[1].Name, "liubin")
	assert.Equal(t, c[1].Password, "qingning")
}

func TestDeleteOne(t *testing.T) {
	prepareTestDatabase()
	c := &CodeBook{ID: 7}
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Delete(c)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(1))
}

func TestDeleteMany(t *testing.T) {
	prepareTestDatabase()
	c := make([]CodeBook, 0)
	c = append(c, CodeBook{ID: 7}, CodeBook{ID: 1}, CodeBook{ID: 8})
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Delete(c)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(2))
}

func TestInsertOne(t *testing.T) {
	prepareTestDatabase()
	c := &CodeBook{Name: "lufei", Password: "lufei", Remarks: nil}
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Insert(c)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(1))
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "lufei"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "lufei")
}

func TestInsertMany(t *testing.T) {
	prepareTestDatabase()
	cc := make([]*CodeBook, 0)
	cc = append(cc, &CodeBook{Name: "xiangjishi", Password: "xiangjishi"})
	cc = append(cc, &CodeBook{Name: "suolong", Password: "suolong"})
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Insert(&cc)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(2))
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "xiangjishi"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "xiangjishi")
}
