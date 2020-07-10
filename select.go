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
	table       Builder
	cols        []string
	distinct    []string
	joins       joins
	wherePreds  predicates
	havingPreds predicates
	limit       *int
	offset      *int
	groupBys    []string
	orderBys    []string
	rebinder    Rebinder
}

func Select(cols ...string) *selectQuery {
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	return &selectQuery{
		cols:     cols,
		groupBys: make([]string, 0),
		orderBys: make([]string, 0),
	}
}

func (q *selectQuery) Select(cols ...string) *selectQuery {
	q.cols = append(q.cols, cols...)
	return q
}

func (q *selectQuery) Distinct(cols ...string) *selectQuery {
	if q.distinct == nil {
		q.distinct = make([]string, 0)
	}
	q.distinct = append(q.distinct, cols...)
	return q
}

func (q *selectQuery) SetCols(cols ...string) *selectQuery {
	q.cols = cols
	return q
}

func (q *selectQuery) From(table Builder) *selectQuery {
	q.table = table
	return q
}

func (q *selectQuery) InnerJoin(table string, condition Builder) *selectQuery {
	q.joins = append(q.joins, newJoin(innerJoin, table, condition))
	return q
}

func (q *selectQuery) LeftJoin(table string, condition Builder) *selectQuery {
	q.joins = append(q.joins, newJoin(leftOuterJoin, table, condition))
	return q
}

func (q *selectQuery) RightJoin(table string, condition Builder) *selectQuery {
	q.joins = append(q.joins, newJoin(rightOuterJoin, table, condition))
	return q
}

func (q *selectQuery) FullJoin(table string, condition Builder) *selectQuery {
	q.joins = append(q.joins, newJoin(fullOuterJoin, table, condition))
	return q
}

func (q *selectQuery) CrossJoin(table string, condition Builder) *selectQuery {
	q.joins = append(q.joins, newJoin(crossJoin, table, condition))
	return q
}

func (q *selectQuery) Where(pred Builder) *selectQuery {
	q.wherePreds = append(q.wherePreds, pred)
	return q
}

func (q *selectQuery) Limit(l int) *selectQuery {
	q.limit = &l
	return q
}

func (q *selectQuery) ClearLimit() *selectQuery {
	q.limit = nil
	return q
}

func (q *selectQuery) Offset(o int) *selectQuery {
	q.offset = &o
	return q
}

func (q *selectQuery) ClearOffset() *selectQuery {
	q.offset = nil
	return q
}

func (q *selectQuery) GroupBy(cols ...string) *selectQuery {
	q.groupBys = append(q.groupBys, cols...)
	return q
}

func (q *selectQuery) Having(pred Builder) *selectQuery {
	q.havingPreds = append(q.havingPreds, pred)
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
	if q.table == nil {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}
	var where string
	var err error

	sb.WriteString("SELECT ")
	if len(q.distinct) > 0 {
		fmt.Fprintf(&sb, "DISTINCT ON (%s) ", strings.Join(q.distinct, ", "))
	} else if q.distinct != nil {
		sb.WriteString("DISTINCT ")
	}
	sb.WriteString(strings.Join(q.cols, ", "))

	switch v := q.table.(type) {
	case S:
		s, _, _ := v.Build()
		fmt.Fprintf(&sb, " FROM %s", s)
	case Builder:
		s, p, err := v.Build()
		if err != nil {
			return "", nil, err
		}
		params = append(params, p...)
		fmt.Fprintf(&sb, " FROM (%s) AS t", s)
	default:
		return "", nil, ErrInvalidTable
	}

	if len(q.joins) > 0 {
		j, _, _ := q.joins.Build()
		fmt.Fprintf(&sb, " %s", j)
	}

	if len(q.wherePreds) > 0 {
		where, params, err = q.wherePreds.Build()
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	if len(q.groupBys) > 0 {
		fmt.Fprintf(&sb, " GROUP BY %s", strings.Join(q.groupBys, ", "))
	}

	if len(q.havingPreds) > 0 {
		having, p, err := q.havingPreds.Build()
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
		query = q.rebinder.Rebind(query)
	}

	return query, params, nil
}
