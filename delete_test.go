package qb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDelete(t *testing.T) {
	q := DeleteFrom("test_table").Where("a", "=", "b").OrWhere("c", ">", 5)

	query, params, err := q.String()
	queryWant := "DELETE FROM test_table WHERE a=? OR c>?"
	paramsWant := []interface{}{"b", 5}

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
