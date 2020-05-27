package qb

import (
	"fmt"
	"strings"
)

type insertQuery struct {
	table     string
	valMap    map[string]interface{}
	returning []string
	err       error
	*conflictResolver
	rebinder Rebinder
}

func InsertInto(table string) *insertQuery {
	return &insertQuery{table: table, valMap: make(map[string]interface{})}
}

func (q *insertQuery) Col(col string, val interface{}) *insertQuery {
	q.valMap[col] = val
	return q
}

func (q *insertQuery) Cols(cols []string, vals ...interface{}) *insertQuery {
	if len(cols) != len(vals) {
		q.err = ErrColValMismatch
		return q
	}

	for i, c := range cols {
		q.Col(c, vals[i])
	}

	return q
}

// OnConflict adds a `ON CONFLICT target action` clause to the query.
// This is only for PostgreSQL.
//
// The target can take three forms:
//		1. (column_name) - a column name (TargetColumn)
//		2. ON CONSTRAINT constraint_name - a UNIQUE constraint name (targetConstraint)
//		3. WHERE predicate - a where clause (whereClause)
//
// The target can take two forms:
// 		1. DO NOTHING - nothing is done (ActionDoNothing)
//		2. DO UPDATE SET col_1=val_1,... WHERE condition - update fields (updateQuery)
func (q *insertQuery) OnConflict(target, action interface{}) *insertQuery {
	q.conflictResolver = &conflictResolver{target, action}
	return q
}

func (q *insertQuery) Returning(cols ...string) *insertQuery {
	q.returning = append(q.returning, cols...)
	return q
}

func (q *insertQuery) RebindWith(r Rebinder) *insertQuery {
	q.rebinder = r
	return q
}

func (q *insertQuery) Build() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	} else if q.err != nil {
		return "", nil, q.err
	}

	keys := orderKeys(q.valMap)
	vals := make([]interface{}, len(keys))
	for i, k := range keys {
		vals[i] = q.valMap[k]
	}

	query := fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES %s`,
		q.table,
		strings.Join(keys, ", "),
		GeneratePlaceholders("?", len(vals)),
	)

	params := vals

	if q.conflictResolver != nil {
		cQuery, p, err := q.conflictResolver.Build()
		if err != nil {
			return "", nil, err
		}
		query = fmt.Sprintf("%s %s", query, cQuery)
		params = append(params, p...)
	}

	if len(q.returning) > 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(q.returning, ", "))
	}

	if q.rebinder != nil {
		query = q.rebinder(query)
	}

	return query, params, nil
}

func (q *insertQuery) String() string {
	query, _, _ := q.Build()
	return query
}
