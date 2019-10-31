package qb

import (
	"fmt"
	"strings"
	"sync"
)

type OrderDir string

const (
	Asc  OrderDir = "ASC"
	Desc OrderDir = "DESC"
)

// selectQueryTS represents a thread safe version of selectQuery.
type selectQueryTS struct {
	mutex sync.RWMutex
	query *selectQuery
}

type selectQuery struct {
	table string
	cols  []string
	joinClause
	whereClause
	limit    *int
	offset   *int
	groupBys []string
	orderBys []string
	rebinder Rebinder
}

func SelectTS(vals ...string) *selectQueryTS {
	if len(vals) == 0 {
		vals = []string{"*"}
	}
	return &selectQueryTS{
		mutex: sync.RWMutex{},
		query: &selectQuery{
			cols:     vals,
			groupBys: make([]string, 0),
			orderBys: make([]string, 0),
		},
	}
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

func (q *selectQueryTS) From(val string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.From(val)
	return q
}
func (q *selectQuery) From(val string) *selectQuery {
	q.table = val
	return q
}

func (q *selectQueryTS) InnerJoin(table, condition string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.InnerJoin(table, condition)
	return q
}
func (q *selectQuery) InnerJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(innerJoin, table, condition))
	return q
}

func (q *selectQueryTS) LeftJoin(table, condition string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.LeftJoin(table, condition)
	return q
}
func (q *selectQuery) LeftJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(leftOuterJoin, table, condition))
	return q
}

func (q *selectQueryTS) RightJoin(table, condition string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.RightJoin(table, condition)
	return q
}
func (q *selectQuery) RightJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(rightOuterJoin, table, condition))
	return q
}

func (q *selectQueryTS) FullJoin(table, condition string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.FullJoin(table, condition)
	return q
}
func (q *selectQuery) FullJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(fullOuterJoin, table, condition))
	return q
}

func (q *selectQueryTS) CrossJoin(table, condition string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.CrossJoin(table, condition)
	return q
}
func (q *selectQuery) CrossJoin(table, condition string) *selectQuery {
	q.joinClause = append(q.joinClause, newJoin(crossJoin, table, condition))
	return q
}

func (q *selectQueryTS) Where(col, cmp string, val interface{}) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.Where(col, cmp, val)
	return q
}
func (q *selectQuery) Where(col, cmp string, val interface{}) *selectQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereAnd})
	return q
}

func (q *selectQueryTS) OrWhere(col, cmp string, val interface{}) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.OrWhere(col, cmp, val)
	return q
}
func (q *selectQuery) OrWhere(col, cmp string, val interface{}) *selectQuery {
	q.clauses = append(q.clauses, clause{col: col, cmp: cmp, val: val, link: whereOr})
	return q
}

func (q *selectQueryTS) Limit(val int) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.Limit(val)
	return q
}
func (q *selectQuery) Limit(val int) *selectQuery {
	q.limit = &val
	return q
}

func (q *selectQueryTS) Offset(val int) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.Offset(val)
	return q
}
func (q *selectQuery) Offset(val int) *selectQuery {
	q.offset = &val
	return q
}

func (q *selectQueryTS) GroupBy(vals ...string) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.GroupBy(vals...)
	return q
}
func (q *selectQuery) GroupBy(vals ...string) *selectQuery {
	q.groupBys = append(q.groupBys, vals...)
	return q
}

func (q *selectQueryTS) OrderBy(col string, dir OrderDir) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.OrderBy(col, dir)
	return q
}
func (q *selectQuery) OrderBy(col string, dir OrderDir) *selectQuery {
	q.orderBys = append(q.orderBys, fmt.Sprintf("%s %s", col, dir))
	return q
}

func (q *selectQueryTS) RebindWith(r Rebinder) *selectQueryTS {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.query.RebindWith(r)
	return q
}
func (q *selectQuery) RebindWith(r Rebinder) *selectQuery {
	q.rebinder = r
	return q
}

func (q *selectQueryTS) String() (string, []interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.query.String()
}
func (q *selectQuery) String() (string, []interface{}, error) {
	if q.table == "" {
		return "", nil, ErrMissingTable
	}

	var sb strings.Builder
	var params []interface{}
	var where string

	fmt.Fprintf(&sb, "SELECT %s FROM %s", strings.Join(q.cols, ", "), q.table)

	if len(q.joinClause) > 0 {
		j, _, _ := q.joinClause.String()
		fmt.Fprintf(&sb, " %s", j)
	}

	if len(q.clauses) > 0 {
		where, params = q.whereClause.string()
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	if len(q.groupBys) > 0 {
		fmt.Fprintf(&sb, " GROUP BY %s", strings.Join(q.groupBys, ", "))
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
