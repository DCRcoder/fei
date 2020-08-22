
# Fei(胐胐)
注: 又北四十里，曰霍山，其木多榖。有兽焉，其状如狸，而白尾，有鬣，名曰腓腓，养之可以已忧。——《山海经·卷五·中山经》

mysql experimental ORM just for fun. it base on [squirrel](https://github.com/Masterminds/squirrel)

### Init
```go
import (
    	_ "github.com/go-sql-driver/mysql"
)

func main() {
    // init db
    engine, err := NewEngine("mysql", dbAddr)
    // init db 读写分离
    engine, err := NewEngineWithMS("mysql", masterAddr, []slaveAddr)
    // init db with config
    // type Config struct {
	//    Driver       string
	//    MasterAddr   string
	//    SlavesAddr   []string
	//    MaxIdleConns int
	//    MaxOpenConns int
	//    Logger       Logger
	//    LogLevel     LogLevel
    // }
    engine, err := New(cfg)
}

```

### Session
```go
    session := engine.NewSession()
    // with context
    session := engine.NewSessionCtx(ctx)
```

### Select
```go
    // count
    engine.NewSession().Select().From("tableName").Where(Eq{"field": "someting"}).Count()
    // findOne
    user := &User{} 
    engine.NewSession().Select().From("tableName").Where(Eq{"field": "someting"}).FindOne(user) // need pointer
    // fineAll
    users := make([]*User, 0)
    engine.NewSession().Select().From("tableName").Where(Eq{"field": "someting"}).FindAll(&user) // need pointer
```

### Insert
```go
    // insert a record
    user := &User{}
    engine.NewSession().Insert(user)
    // insert multiple
    users := []*User{}
    engine.NewSession().Insert(users)
```

### Update
```go
    // update a record
    user := &User{}
    engine.NewSession().Update(user)
    // update multiple
    users := []*User{}
    engine.NewSession().Update(users)

```

### Delete
```go
    // delete a record
    user := &User{}
    engine.NewSession().Delete(user)
    // delete multiple
    users := []*User{}
    engine.NewSession().Delete(users)
```

### sql driver also support
```go
    // you can also use sql Query, Exec
    engine.NewSession().Query("select * from user")
    engine.NewSession().QueryContext(ctx, "select * from user")
    engine.NewSession().Exec("insert ....")
    engine.NewSession().ExecContext("insert ....")
```

### TODO
- [ ] Transaction
- [ ] plugin
