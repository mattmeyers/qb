package qb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInsert(t *testing.T) {
	q := InsertInto("test_table").Columns("a", "b").Values("c", 1).Columns("d").Values(false)

	query, params, err := q.String()
	queryWant := "INSERT INTO test_table (a, b, d) VALUES (?, ?, ?)"
	paramsWant := []interface{}{"c", 1, false}

	if err != nil {
		t.Errorf("String() failed: %s", err)
	}

	if query != queryWant {
		t.Errorf("Wrong query: want %s, got %s", queryWant, query)
	}
	if !cmp.Equal(params, paramsWant) {
		t.Errorf("Wrong params: want %v, got %v", paramsWant, params)
	}
}
