package qb

import (
	"fmt"
	"strings"
)

type updateQuery struct {
	table    string
	setPairs map[string]interface{}
	whereClause
}

func Update(table string) *updateQuery {
	return &updateQuery{table: table, setPairs: make(map[string]interface{})}
}

func (q *updateQuery) Set(col string, val interface{}) *updateQuery {
	q.setPairs[col] = val
	return q
}

func (q *updateQuery) Where(col, cmp string, val interface{}) *updateQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *updateQuery) OrWhere(col, cmp string, val interface{}) *updateQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *updateQuery) String() (string, []interface{}, error) {
	return q.string(true)
}

func (q *updateQuery) string(tableRequired bool) (string, []interface{}, error) {
	if q.table == "" && tableRequired {
		return "", nil, ErrMissingTable
	} else if len(q.setPairs) == 0 {
		return "", nil, ErrMissingSetPairs
	}

	var sb strings.Builder
	params := make([]interface{}, len(q.setPairs))

	sb.WriteString("UPDATE ")

	if q.table != "" {
		sb.WriteString(q.table)
		sb.WriteString(" ")
	}
	sb.WriteString("SET ")

	keys := orderKeys(q.setPairs)

	sets := make([]string, len(q.setPairs))
	for i, k := range keys {
		sets[i] = fmt.Sprintf("%s=?", k)
		params[i] = q.setPairs[k]
		i++
	}
	sb.WriteString(strings.Join(sets, ", "))

	if len(q.clauses) > 0 {
		sb.WriteString(" WHERE ")
		where, wParams := q.whereClause.string()
		params = append(params, wParams...)
		sb.WriteString(where)
	}

	return sb.String(), params, nil
}
