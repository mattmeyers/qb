package qb

import (
	"fmt"
	"strings"
)

type Excluded string

type updateQuery struct {
	table      string
	setPairs   map[string]interface{}
	wherePreds predicates
	rebinder   Rebinder
}

func Update(table string) *updateQuery {
	return &updateQuery{table: table, setPairs: make(map[string]interface{})}
}

func (q *updateQuery) Set(col string, val interface{}) *updateQuery {
	q.setPairs[col] = val
	return q
}

func (q *updateQuery) Where(pred Builder) *updateQuery {
	q.wherePreds = append(q.wherePreds, pred)
	return q
}

func (q *updateQuery) RebindWith(r Rebinder) *updateQuery {
	q.rebinder = r
	return q
}

func (q *updateQuery) Build() (string, []interface{}, error) {
	return q.build(true)
}

func (q *updateQuery) build(tableRequired bool) (string, []interface{}, error) {
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

	if len(q.wherePreds) > 0 {
		q, p, err := q.wherePreds.Build()
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(" WHERE ")
		sb.WriteString(q)
		params = append(params, p...)
	}

	query := sb.String()

	if q.rebinder != nil {
		query = q.rebinder.Rebind(query)
	}

	return query, params, nil
}
