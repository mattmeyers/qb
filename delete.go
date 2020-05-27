package qb

import (
	"strings"
)

type deleteQuery struct {
	table        string
	whereClauses whereClause
	returning    []string
}

func DeleteFrom(table string) *deleteQuery {
	return &deleteQuery{table: table}
}

func (q *deleteQuery) Where(clause Builder) *deleteQuery {
	q.whereClauses.clauses = append(q.whereClauses.clauses, clause)
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

	if len(q.whereClauses.clauses) > 0 {
		w, p, err := q.whereClauses.Build()
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
