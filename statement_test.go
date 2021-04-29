package fei

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
)

func TestSelectStatement(t *testing.T) {
	st := &Statement{}
	st.Select("*").From("user").Where(Eq{"name": "laojun", "sex": []string{"male"}, "age": nil}).OrderBy("id")
	s, _, err := st.ToSQL()
	assert.Equal(t, err, nil)
	fmt.Println(s)
	assert.Equal(t, s, "SELECT * FROM user WHERE age IS NULL AND name = ? AND sex IN (?) ORDER BY id")

	or := make(OR, 0)
	or = append(or, Eq{"name": []string{"qingning"}}, GT{"age": 1000})
	st.Select("*").From("user").Where(or)
	s, _, err = st.ToSQL()
	assert.Equal(t, err, nil)
	fmt.Println(s)
	assert.Equal(t, s, "SELECT * FROM user WHERE (name IN (?) OR age > ?)")

	st.Select("*").From("user").EnableExplain(true)
	s, _, err = st.ToSQL()
	assert.Equal(t, err, nil)
	fmt.Println(s)
	assert.Equal(t, s, "EXPLAIN SELECT * FROM user")

	st.Select("*").From("user").UseIndexs("idx", "idxy")
	s, _, err = st.ToSQL()
	assert.Equal(t, err, nil)
	fmt.Println(s)
	assert.Equal(t, s, "SELECT * FROM user USE INDEX (idx, idxy)")
}
