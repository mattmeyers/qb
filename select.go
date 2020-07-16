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

type SelectQuery struct {
	table        string
	fromSub      *SelectQuery
	fromSubAlias string
	cols         []string
	distinct     []string
	joins        joins
	wherePreds   predicates
	havingPreds  predicates
	limit        *int
	offset       *int
	groupBys     []string
	orderBys     []string
	rebinder     Rebinder
}

func Select(cols ...string) *SelectQuery {
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	return &SelectQuery{
		cols:     cols,
		groupBys: make([]string, 0),
		orderBys: make([]string, 0),
	}
}

func (q *SelectQuery) Select(cols ...string) *SelectQuery {
	q.cols = append(q.cols, cols...)
	return q
}

func (q *SelectQuery) Distinct(cols ...string) *SelectQuery {
	if q.distinct == nil {
		q.distinct = make([]string, 0)
	}
	q.distinct = append(q.distinct, cols...)
	return q
}

func (q *SelectQuery) SetCols(cols ...string) *SelectQuery {
	q.cols = cols
	return q
}

func (q *SelectQuery) From(table string) *SelectQuery {
	q.table = table
	return q
}

func (q *SelectQuery) FromSub(query *SelectQuery, alias string) *SelectQuery {
	q.fromSub = query

	if alias == "" {
		alias = "from_sub"
	}
	q.fromSubAlias = alias

	return q
}

func (q *SelectQuery) InnerJoin(table string, condition Builder) *SelectQuery {
	q.joins = append(q.joins, newJoin(innerJoin, table, condition))
	return q
}

func (q *SelectQuery) LeftJoin(table string, condition Builder) *SelectQuery {
	q.joins = append(q.joins, newJoin(leftOuterJoin, table, condition))
	return q
}

func (q *SelectQuery) RightJoin(table string, condition Builder) *SelectQuery {
	q.joins = append(q.joins, newJoin(rightOuterJoin, table, condition))
	return q
}

func (q *SelectQuery) FullJoin(table string, condition Builder) *SelectQuery {
	q.joins = append(q.joins, newJoin(fullOuterJoin, table, condition))
	return q
}

func (q *SelectQuery) CrossJoin(table string, condition Builder) *SelectQuery {
	q.joins = append(q.joins, newJoin(crossJoin, table, condition))
	return q
}

func (q *SelectQuery) Where(pred Builder) *SelectQuery {
	q.wherePreds = append(q.wherePreds, pred)
	return q
}

func (q *SelectQuery) Limit(l int) *SelectQuery {
	q.limit = &l
	return q
}

func (q *SelectQuery) ClearLimit() *SelectQuery {
	q.limit = nil
	return q
}

func (q *SelectQuery) Offset(o int) *SelectQuery {
	q.offset = &o
	return q
}

func (q *SelectQuery) ClearOffset() *SelectQuery {
	q.offset = nil
	return q
}

func (q *SelectQuery) GroupBy(cols ...string) *SelectQuery {
	q.groupBys = append(q.groupBys, cols...)
	return q
}

func (q *SelectQuery) Having(pred Builder) *SelectQuery {
	q.havingPreds = append(q.havingPreds, pred)
	return q
}

func (q *SelectQuery) OrderBy(col string, dir OrderDir) *SelectQuery {
	q.orderBys = append(q.orderBys, fmt.Sprintf("%s %s", col, dir))
	return q
}

func (q *SelectQuery) RebindWith(r Rebinder) *SelectQuery {
	q.rebinder = r
	return q
}

func (q *SelectQuery) String() string {
	s, _, _ := q.Build()
	return s
}

func (q *SelectQuery) Build() (string, []interface{}, error) {
	if (q.table == "" && q.fromSub == nil) || (q.table != "" && q.fromSub != nil) {
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

	if q.table != "" {
		fmt.Fprintf(&sb, " FROM %s", q.table)
	} else {
		s, p, err := q.fromSub.Build()
		if err != nil {
			return "", nil, err
		}
		params = append(params, p...)
		fmt.Fprintf(&sb, " FROM (%s) AS t", s)
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
