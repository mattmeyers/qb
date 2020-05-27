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

func (jt joinType) toString() string {
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
	condition string
}

type joinClause []join

func newJoin(joinType joinType, table, condition string) join {
	return join{joinType, table, condition}
}

func (jc joinClause) Build() (string, []interface{}, error) {
	parts := make([]string, len(jc))
	for i, j := range jc {
		parts[i] = fmt.Sprintf("%s %s ON %s", j.joinType.toString(), j.table, j.condition)
	}
	return strings.Join(parts, " "), nil, nil
}
