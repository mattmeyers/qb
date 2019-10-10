package qb

import (
	"fmt"
	"strings"
)

type queryType int

const (
	querySelect queryType = iota
	queryInsert
	queryUpdate
	queryDelete
)

type Query struct {
	queryType    queryType
	table        string
	selectFields []string
	insertCols   []string
	values       []interface{}
	where        whereClauses
}

type whereLink int

const (
	whereStart whereLink = iota
	whereAnd
	whereOr
)

type whereClause struct {
	col  string
	cmp  string
	val  interface{}
	link whereLink
}
type whereClauses []whereClause

func New() *Query {
	return &Query{}
}

func (q *Query) Select(vals ...string) *Query {
	q.queryType = querySelect
	q.selectFields = append(q.selectFields, vals...)
	return q
}

func (q *Query) From(val string) *Query {
	q.table = val
	return q
}

func (q *Query) InsertInto(val string) *Query {
	q.queryType = queryInsert
	q.table = val
	return q
}

func (q *Query) Columns(vals ...string) *Query {
	q.insertCols = append(q.insertCols, vals...)
	return q
}

func (q *Query) Values(vals ...interface{}) *Query {
	q.values = append(q.values, vals...)
	return q
}

func (q *Query) Where(col, cmp string, val interface{}) *Query {
	q.where = append(q.where, whereClause{col: col, cmp: cmp, val: val, link: whereStart})
	return q
}

func (q *Query) AndWhere(col, cmp string, val interface{}) *Query {
	q.where = append(q.where, whereClause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *Query) OrWhere(col, cmp string, val interface{}) *Query {
	q.where = append(q.where, whereClause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *Query) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var query string
	var params []interface{}

	switch q.queryType {
	case querySelect:
		query, params = q.buildSelect()
	case queryInsert:
		query, params = q.buildInsert()
	}

	return query, params, nil
}

func (q *Query) buildSelect() (string, []interface{}) {
	var sb strings.Builder
	var params []interface{}
	var where string

	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.selectFields, ", "), q.table))

	if len(q.where) > 0 {
		where, params = q.where.string()
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	return sb.String(), params
}

func (q *Query) buildInsert() (string, []interface{}) {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, strings.Join(q.insertCols, ", "), strings.Join(strings.Split(strings.Repeat("?", len(q.values)), ""), ", "))

	return query, q.values
}

func (q *Query) buildUpdate() (string, []interface{}) {
	return "", nil
}

func (q *Query) buildDelete() (string, []interface{}) {
	return "", nil
}

func (w whereClauses) string() (string, []interface{}) {
	var parts []string
	var params []interface{}

	for i, c := range w {
		if i > 0 {
			if c.link == whereAnd {
				parts = append(parts, "AND")
			} else if c.link == whereOr {
				parts = append(parts, "OR")
			}
		}

		parts = append(parts, fmt.Sprintf("%s%s?", c.col, c.cmp))
		params = append(params, c.val)

	}
	return strings.Join(parts, " "), params
}
