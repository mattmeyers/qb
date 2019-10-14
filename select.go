package qb

import (
	"fmt"
	"strings"
)

type selectQuery struct {
	table string
	cols  []string
	whereClause
}

func Select(vals ...string) *selectQuery {
	if len(vals) == 0 {
		vals = []string{"*"}
	}
	return &selectQuery{cols: vals}
}

func (q *selectQuery) From(val string) *selectQuery {
	q.table = val
	return q
}

func (q *selectQuery) Where(col, cmp string, val interface{}) *selectQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *selectQuery) OrWhere(col, cmp string, val interface{}) *selectQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *selectQuery) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}
	var where string

	fmt.Fprintf(&sb, "SELECT %s FROM %s", strings.Join(q.cols, ", "), q.table)

	if len(q.clauses) > 0 {
		where, params = q.whereClause.string()
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	return sb.String(), params, nil
}
