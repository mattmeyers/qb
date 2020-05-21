package qb

import (
	"sort"
	"strings"
)

type QueryBuilder interface {
	SQL() (string, []interface{}, error)
}

type raw struct {
	q string
	p []interface{}
}

func Raw(q string, p []interface{}) raw { return raw{q: q, p: p} }

func (r raw) SQL() (string, []interface{}, error) { return r.q, r.p, nil }

// Rebinder represents a function that can replace all `?` tokens in the query
// with dialect specific tokens. These other tokens are dialect and driver
// specific.
type Rebinder func(string) string

// GeneratePlaceholders generates a comma seperated list of the provided
// symbol and places the list in parentheses. If num is less than or
// equal to zero, then an empty set of parentheses is returned.
func GeneratePlaceholders(symbol string, num int) string {
	if num <= 0 {
		return "()"
	}

	p := make([]string, num)
	for i := 0; i < num; i++ {
		p[i] = symbol
	}
	return "(" + strings.Join(p, ", ") + ")"
}

func orderKeys(m map[string]interface{}) []string {
	kArr := make([]string, len(m))
	i := 0
	for k := range m {
		kArr[i] = k
		i++
	}
	sort.Strings(kArr)
	return kArr
}
