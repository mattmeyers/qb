package qb

import (
	"sort"
	"strings"
)

type Builder interface {
	Build() (string, []interface{}, error)
}

type S string

func (s S) Build() (string, []interface{}, error) { return string(s), nil, nil }

type raw struct {
	q string
	p []interface{}
}

func Raw(q string, p []interface{}) raw { return raw{q: q, p: p} }

func (r raw) Build() (string, []interface{}, error) { return r.q, r.p, nil }

// Rebinder represents a function that can replace all `?` tokens in the query
// with dialect specific tokens. These other tokens are dialect and driver
// specific.
type Rebinder interface {
	Rebind(string) string
}

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
