package qb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSelect(t *testing.T) {

	q := Select("a", "b", "c").From("test_table").Where("d", "=", "e").Where("f", ">", 5).OrWhere("g", "!=", false)

	query, params, err := q.String()
	queryWant := "SELECT a, b, c FROM test_table WHERE d=? AND f>? OR g!=?"
	paramsWant := []interface{}{"e", 5, false}

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
