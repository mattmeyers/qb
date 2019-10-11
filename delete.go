package qb

import "strings"

type deleteQuery struct {
	table string
	whereClause
}

func DeleteFrom(table string) *deleteQuery {
	return &deleteQuery{table: table}
}

func (q *deleteQuery) Where(col, cmp string, val interface{}) *deleteQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *deleteQuery) OrWhere(col, cmp string, val interface{}) *deleteQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *deleteQuery) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}
	var where string

	sb.WriteString("DELETE FROM ")
	sb.WriteString(q.table)

	if len(q.clauses) > 0 {
		sb.WriteString(" WHERE ")
		where, params = q.whereClause.string()
		sb.WriteString(where)
	}

	return sb.String(), params, nil
}
