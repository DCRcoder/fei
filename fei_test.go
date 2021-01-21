package fei

import (
	"database/sql"
	"fmt"
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
	ID        int64 `fei:"pk"`
	Name      string
	Password  string
	Remarks   *string
	CreatedAt time.Time `fei:"readOnly"`
	UpdatedAt string    `fei:"readOnly"`
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
	assert.Equal(t, c.ID, int64(7))
	assert.Equal(t, c.Name, "laojun")
	assert.Equal(t, *c.Remarks, "qingning")
	assert.NotEqual(t, c.CreatedAt, "")
	assert.NotEqual(t, c.UpdatedAt, nil)
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
	assert.Equal(t, c[1].ID, int64(9))
	assert.Equal(t, c[1].Name, "liubin")
	assert.Equal(t, c[1].Password, "qingning")
	cc := make([]*CodeBook, 0)
	err = engine.NewSession().Select().FindAll(&cc)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(cc), 5)
	assert.Equal(t, cc[1].ID, int64(7))
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

func TestUpdateOne(t *testing.T) {
	prepareTestDatabase()
	c := &CodeBook{ID: 2, Name: "nami", Password: "lufei", Remarks: nil}
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Update(c)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(1))
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "nami"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "lufei")
}

func TestUpdateMany(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	cs := []*CodeBook{
		&CodeBook{ID: 2, Name: "nami", Password: "lufei", Remarks: nil},
		&CodeBook{ID: 10, Name: "nami123", Password: "feifei", Remarks: nil},
	}
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().Update(cs)
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(2))
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "nami123"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "feifei")
}

func TestUpdateRowOne(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().From("codebook").Where(Eq{"id": 2}).UpdateRow(map[string]interface{}{"name": "maomao"})
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(1))
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"id": 2}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Name, "maomao")
}

func TestUpdateRowMany(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	rowcount, err := engine.NewSession().From("codebook").Where(Eq{"id": []int{2, 10}}).UpdateRow(map[string]interface{}{"name": "maomao"})
	assert.Equal(t, err, nil)
	assert.Equal(t, rowcount, int64(2))
	nc := make([]*CodeBook, 0)
	err = engine.NewSession().Select().Where(Eq{"id": []int{2, 10}}).FindAll(&nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc[0].Name, "maomao")
	assert.Equal(t, nc[1].Name, "maomao")
}

func TestTransaction(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	session := engine.NewSession()
	f := func(s *Session) (interface{}, error) {
		c := &CodeBook{ID: 2, Name: "nami", Password: "lufei", Remarks: nil}
		_, err := s.Update(c)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err = session.Transaction(f)
	assert.Equal(t, err, nil)
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "nami"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "lufei")
}

func TestTransactionWithError(t *testing.T) {
	prepareTestDatabase()
	engine, err := NewEngine("mysql", dbAddr)
	engine.SetLogLevel(LogDebug)
	assert.Equal(t, err, nil)
	session := engine.NewSession()
	f := func(s *Session) (interface{}, error) {
		c := &CodeBook{ID: 2, Name: "nami", Password: "lufei", Remarks: nil}
		_, err := s.Update(c)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("update error")
	}
	_, err = session.Transaction(f)
	assert.NotEqual(t, err, nil)
	nc := &CodeBook{}
	err = engine.NewSession().Select().Where(Eq{"name": "nami"}).FindOne(nc)
	assert.Equal(t, err, nil)
	assert.Equal(t, nc.Password, "nami")
}
