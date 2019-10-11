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

type clause struct {
	col  string
	cmp  string
	val  interface{}
	link whereLink
}

type whereClause struct {
	clauses []clause
}

func (w whereClause) string() (string, []interface{}) {
	var parts []string
	var params []interface{}

	for i, c := range w.clauses {
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
