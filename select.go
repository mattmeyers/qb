package qb

import (
	"fmt"
	"strings"
)

type OrderDir string

const (
	Asc  OrderDir = "ASC"
	Desc OrderDir = "DESC"
)

type selectQuery struct {
	table interface{}
	cols  []string
	joinClause
	whereClauses  whereClause
	havingClauses whereClause
	limit         *int
	offset        *int
	groupBys      []string
	orderBys      []string
	rebinder      Rebinder
}

func Select(vals ...string) *selectQuery {
	if len(vals) == 0 {
		vals = []string{"*"}
	}
	return &selectQuery{
		cols:     vals,
		groupBys: make([]string, 0),
		orderBys: make([]string, 0),
	}
}

func (q *selectQuery) Select(vals ...string) *selectQuery {
	q.cols = append(q.cols, vals...)
	return q
}

func (q *selectQuery) SetCols(vals ...string) *selectQuery {
	q.cols = vals
	return q
}

func (q *selectQuery) From(val interface{}) *selectQuery {
	q.table = val
	return q
}

func (q *selectQuery) InnerJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(innerJoin, table, condition))
	return q
}

func (q *selectQuery) LeftJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(leftOuterJoin, table, condition))
	return q
}

func (q *selectQuery) RightJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(rightOuterJoin, table, condition))
	return q
}

func (q *selectQuery) FullJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(fullOuterJoin, table, condition))
	return q
}

func (q *selectQuery) CrossJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(crossJoin, table, condition))
	return q
}

func (q *selectQuery) Where(clause Builder) *selectQuery {
	q.whereClauses.clauses = append(q.whereClauses.clauses, clause)
	return q
}

func (q *selectQuery) Limit(val int) *selectQuery {
	q.limit = &val
	return q
}

func (q *selectQuery) ClearLimit() *selectQuery {
	q.limit = nil
	return q
}

func (q *selectQuery) Offset(val int) *selectQuery {
	q.offset = &val
	return q
}

func (q *selectQuery) ClearOffset() *selectQuery {
	q.offset = nil
	return q
}

func (q *selectQuery) GroupBy(vals ...string) *selectQuery {
	q.groupBys = append(q.groupBys, vals...)
	return q
}

func (q *selectQuery) Having(clause Builder) *selectQuery {
	q.havingClauses.clauses = append(q.havingClauses.clauses, clause)
	return q
}

func (q *selectQuery) OrderBy(col string, dir OrderDir) *selectQuery {
	q.orderBys = append(q.orderBys, fmt.Sprintf("%s %s", col, dir))
	return q
}

func (q *selectQuery) RebindWith(r Rebinder) *selectQuery {
	q.rebinder = r
	return q
}

func (q *selectQuery) String() string {
	s, _, _ := q.Build()
	return s
}

func (q *selectQuery) Build() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}
	var where string
	var err error

	fmt.Fprintf(&sb, "SELECT %s", strings.Join(q.cols, ", "))

	switch v := q.table.(type) {
	case Builder:
		s, p, err := v.Build()
		if err != nil {
			return "", nil, err
		}
		params = append(params, p...)
		fmt.Fprintf(&sb, " FROM (%s) AS t", s)
	case string:
		fmt.Fprintf(&sb, " FROM %s", v)
	default:
		return "", nil, ErrInvalidTable
	}

	if len(q.joinClause) > 0 {
		j, _, _ := q.joinClause.String()
		fmt.Fprintf(&sb, " %s", j)
	}

	if len(q.whereClauses.clauses) > 0 {
		where, params, err = q.whereClauses.Build()
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	if len(q.groupBys) > 0 {
		fmt.Fprintf(&sb, " GROUP BY %s", strings.Join(q.groupBys, ", "))
	}

	if len(q.havingClauses.clauses) > 0 {
		having, p, err := q.havingClauses.Build()
		if err != nil {
			return "", nil, err
		}

		params = append(params, p...)
		sb.WriteString(" HAVING ")
		sb.WriteString(having)
	}

	if len(q.orderBys) > 0 {
		fmt.Fprintf(&sb, " ORDER BY %s", strings.Join(q.orderBys, ", "))
	}

	if q.limit != nil {
		fmt.Fprintf(&sb, " LIMIT %d", *q.limit)
	}

	if q.offset != nil {
		fmt.Fprintf(&sb, " OFFSET %d", *q.offset)
	}

	query := sb.String()
	if q.rebinder != nil {
		query = q.rebinder(query)
	}

	return query, params, nil
}
