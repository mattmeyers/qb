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
	whereParams  []whereParam
}

type whereParam struct {
	col string
	cmp string
	val interface{}
}

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
	q.whereParams = append(q.whereParams, whereParam{col: col, cmp: cmp, val: val})
	return q
}

func (q *Query) String() (string, []interface{}) {
	var query string
	var params []interface{}

	switch q.queryType {
	case querySelect:
		query, params = q.buildSelect()
	case queryInsert:
		query, params = q.buildInsert()
	}

	return query, params
}

func (q *Query) buildSelect() (string, []interface{}) {
	var sb strings.Builder
	var params []interface{}

	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.selectFields, ", "), q.table))

	if wLen := len(q.whereParams); wLen > 0 {
		sb.WriteString(" WHERE ")

		conditions := make([]string, wLen)
		for i, w := range q.whereParams {
			conditions[i] = fmt.Sprintf("%s%s?", w.col, w.cmp)
			params = append(params, w.val)
		}

		sb.WriteString(strings.Join(conditions, ", "))
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
