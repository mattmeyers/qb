package qb

import (
	"fmt"
	"strings"
)

type Excluded string

type updateQuery struct {
	table    string
	setPairs map[string]interface{}
	whereClause
	rebinder Rebinder
}

func Update(table string) *updateQuery {
	return &updateQuery{table: table, setPairs: make(map[string]interface{})}
}

func (q *updateQuery) Set(col string, val interface{}) *updateQuery {
	q.setPairs[col] = val
	return q
}

func (q *updateQuery) Where(col, cmp string, val interface{}) *updateQuery {
	// q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *updateQuery) OrWhere(col, cmp string, val interface{}) *updateQuery {
	// q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *updateQuery) RebindWith(r Rebinder) *updateQuery {
	q.rebinder = r
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
	params := make([]interface{}, 0)

	sb.WriteString("UPDATE ")

	if q.table != "" {
		fmt.Fprintf(&sb, `"%s" `, q.table)
	}
	sb.WriteString("SET ")

	keys := orderKeys(q.setPairs)

	sets := make([]string, 0)
	for _, k := range keys {
		switch v := q.setPairs[k].(type) {
		case Excluded:
			sets = append(sets, fmt.Sprintf("%s=EXCLUDED.%s", k, v))
		default:
			sets = append(sets, fmt.Sprintf("%s=?", k))
			params = append(params, q.setPairs[k])
		}
	}
	sb.WriteString(strings.Join(sets, ", "))

	// if len(q.clauses) > 0 {
	// 	sb.WriteString(" WHERE ")
	// 	where, wParams := q.whereClause.string()
	// 	params = append(params, wParams...)
	// 	sb.WriteString(where)
	// }

	query := sb.String()

	if q.rebinder != nil {
		query = q.rebinder(query)
	}

	return query, params, nil
}
