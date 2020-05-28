package qb

import (
	"fmt"
	"strings"
)

type joinType int

const (
	innerJoin joinType = iota
	leftOuterJoin
	rightOuterJoin
	fullOuterJoin
	crossJoin
)

func (jt joinType) String() string {
	switch jt {
	case innerJoin:
		return "INNER JOIN"
	case leftOuterJoin:
		return "LEFT OUTER JOIN"
	case rightOuterJoin:
		return "RIGHT OUTER JOIN"
	case fullOuterJoin:
		return "FULL OUTER JOIN"
	case crossJoin:
		return "CROSS JOIN"
	}
	return ""
}

type join struct {
	joinType  joinType
	table     string
	condition interface{}
}

type joins []join

func newJoin(joinType joinType, table string, condition interface{}) join {
	return join{joinType, table, condition}
}

func (jc joins) Build() (string, []interface{}, error) {
	parts := make([]string, len(jc))
	var params []interface{}
	for i, j := range jc {
		switch v := j.condition.(type) {
		case string:
			parts[i] = fmt.Sprintf("%s %s ON %s", j.joinType.String(), j.table, v)
		case Builder:
			q, p, err := v.Build()
			if err != nil {
				return "", nil, err
			}
			parts[i] = fmt.Sprintf("%s %s ON %s", j.joinType.String(), j.table, q)
			params = append(params, p...)
		default:
			return "", nil, ErrInvalidType
		}
	}
	return strings.Join(parts, " "), params, nil
}
