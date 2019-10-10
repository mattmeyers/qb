package qb

import (
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	q := New()
	q.Select("a", "b", "c").From("test_table").Where("d", "=", "e").AndWhere("f", ">", 5).OrWhere("g", "!=", false)

	query, params, _ := q.String()
	fmt.Println(query)
	fmt.Println(params)
}

func TestInsert(t *testing.T) {
	q := New()
	q.InsertInto("test_table").Columns("a", "b").Values("c", 4)

	if true {
		q.Columns("col").Values(true)
	}

	query, params, _ := q.String()
	fmt.Println(query)
	fmt.Println(params)
}
