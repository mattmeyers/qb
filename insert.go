package qb

import (
	"fmt"
	"strings"
)

type insertQuery struct {
	table string
	cols  []string
	vals  []interface{}
}

func InsertInto(table string) *insertQuery {
	return &insertQuery{table: table}
}

func (q *insertQuery) Columns(cols ...string) *insertQuery {
	q.cols = append(q.cols, cols...)
	return q
}

func (q *insertQuery) Values(vals ...interface{}) *insertQuery {
	q.vals = append(q.vals, vals...)
	return q
}

func (q *insertQuery) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		q.table,
		strings.Join(q.cols, ", "),
		GeneratePlaceholders("?", len(q.vals)),
	)

	return query, q.vals, nil
}
