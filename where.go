package qb

import (
	"fmt"
	"strings"
)

type whereLink int

const (
	whereAnd whereLink = iota
	whereOr
)

type Cmp struct {
	Col string
	Op  string
	Val interface{}
}

func Eq(col string, val interface{}) Cmp { return Cmp{Col: col, Op: "=", Val: val} }

func Neq(col string, val interface{}) Cmp { return Cmp{Col: col, Op: "<>", Val: val} }

func Gt(col string, val interface{}) Cmp { return Cmp{Col: col, Op: ">", Val: val} }

func Gte(col string, val interface{}) Cmp { return Cmp{Col: col, Op: ">=", Val: val} }

func Lt(col string, val interface{}) Cmp { return Cmp{Col: col, Op: "<", Val: val} }

func Lte(col string, val interface{}) Cmp { return Cmp{Col: col, Op: "<=", Val: val} }

type Or []QueryBuilder

type And []QueryBuilder

func (o Or) SQL() (string, []interface{}, error) {
	parts := make([]string, len(o))
	params := make([]interface{}, 0, len(o))

	for i, c := range o {
		q, p, err := c.SQL()
		if err != nil {
			return "", nil, err
		}
		parts[i] = q
		params = append(params, p...)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, " OR ")), params, nil
}

func (a And) SQL() (string, []interface{}, error) {
	parts := make([]string, len(a))
	params := make([]interface{}, 0, len(a))

	for i, c := range a {
		q, p, err := c.SQL()
		if err != nil {
			return "", nil, err
		}
		parts[i] = q
		params = append(params, p...)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, " AND ")), params, nil
}

func (c Cmp) SQL() (string, []interface{}, error) {
	var q string
	var p []interface{}
	var err error

	switch v := c.Val.(type) {
	case QueryBuilder:
		q, p, err = v.SQL()
		if err != nil {
			return "", nil, err
		}

		q = fmt.Sprintf("%s%s(%s)", c.Col, c.Op, q)
	default:
		q = fmt.Sprintf("%s%s?", c.Col, c.Op)
		p = []interface{}{c.Val}
	}
	return q, p, nil
}

type whereClause struct {
	clauses []QueryBuilder
}

func (w whereClause) SQL() (string, []interface{}, error) {
	var parts []string
	var params []interface{}

	for _, c := range w.clauses {
		part, param, err := c.SQL()
		if err != nil {
			return "", nil, err
		}

		parts = append(parts, part)
		params = append(params, param...)

	}
	return strings.Join(parts, " AND "), params, nil
}
