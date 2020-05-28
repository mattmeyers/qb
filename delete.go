package qb

import (
	"strings"
)

type deleteQuery struct {
	table      string
	wherePreds predicates
	returning  []string
}

func DeleteFrom(table string) *deleteQuery {
	return &deleteQuery{table: table}
}

func (q *deleteQuery) Where(pred Builder) *deleteQuery {
	q.wherePreds = append(q.wherePreds, pred)
	return q
}

func (q *deleteQuery) Returning(cols ...string) *deleteQuery {
	q.returning = append(q.returning, cols...)
	return q
}

func (q *deleteQuery) Build() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}

	sb.WriteString("DELETE FROM ")
	sb.WriteString(q.table)

	if len(q.wherePreds) > 0 {
		w, p, err := q.wherePreds.Build()
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(" WHERE ")
		sb.WriteString(w)
		params = append(params, p...)
	}

	if len(q.returning) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(strings.Join(q.returning, ", "))
	}

	return sb.String(), params, nil
}

func (q *deleteQuery) String() string {
	query, _, _ := q.Build()
	return query
}
