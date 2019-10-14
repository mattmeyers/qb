package qb

import (
	"fmt"
	"strings"
)

type insertQuery struct {
	table  string
	valMap map[string]interface{}
	err    error
}

func InsertInto(table string) *insertQuery {
	return &insertQuery{table: table, valMap: make(map[string]interface{})}
}

func (q *insertQuery) Col(col string, val interface{}) *insertQuery {
	q.valMap[col] = val
	return q
}

func (q *insertQuery) Cols(cols []string, vals ...interface{}) *insertQuery {
	if len(cols) != len(vals) {
		q.err = ErrColValMismatch
		return q
	}

	for i, c := range cols {
		q.Col(c, vals[i])
	}

	return q
}

func (q *insertQuery) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	} else if q.err != nil {
		return "", nil, q.err
	}

	keys := orderKeys(q.valMap)
	vals := make([]interface{}, len(keys))
	for i, k := range keys {
		vals[i] = q.valMap[k]
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		q.table,
		strings.Join(keys, ", "),
		GeneratePlaceholders("?", len(vals)),
	)

	return query, vals, nil
}
