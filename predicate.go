package qb

import (
	"fmt"
	"strings"
)

// Pred represents a SQL predicate. The convenience methods Eq, Neq, Gt, Gte,
// Lt, and Lte are provided for composing common predicates.
type Pred struct {
	Col string
	Op  string
	Val interface{}
}

// Eq returns a predicate using the `=` operator.
func Eq(col string, val interface{}) Pred { return Pred{Col: col, Op: "=", Val: val} }

// Neq returns a predicate using the `!=` operator.
func Neq(col string, val interface{}) Pred { return Pred{Col: col, Op: "!=", Val: val} }

// Gt returns a predicate using the `>` operator.
func Gt(col string, val interface{}) Pred { return Pred{Col: col, Op: ">", Val: val} }

// Gte returns a predicate using the `>=` operator.
func Gte(col string, val interface{}) Pred { return Pred{Col: col, Op: ">=", Val: val} }

// Lt returns a predicate using the `<` operator.
func Lt(col string, val interface{}) Pred { return Pred{Col: col, Op: "<", Val: val} }

// Lte returns a predicate using the `<=` operator.
func Lte(col string, val interface{}) Pred { return Pred{Col: col, Op: "<=", Val: val} }

// Or implements the Builder interface for a list of Builders. When built, the
// slice of builders are combined with the `OR` operator and the entire
// predicate is surrounded with parentheses.
type Or []Builder

// And implements the Builder interface for a list of Builders. When built, the
// slice of builders are combined with the `AND` operator and the entire
// predicate is surrounded with parentheses.
type And []Builder

// Build creates a predicate by combining the slice of Builders with the `OR`
// operator and surrounding the predicate with parentheses.
func (o Or) Build() (string, []interface{}, error) {
	parts := make([]string, len(o))
	params := make([]interface{}, 0, len(o))

	for i, c := range o {
		q, p, err := c.Build()
		if err != nil {
			return "", nil, err
		}
		parts[i] = q
		params = append(params, p...)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, " OR ")), params, nil
}

// Build creates a predicate by combining the slice of Builders with the `AND`
// operator and surrounding the predicate with parentheses.
func (a And) Build() (string, []interface{}, error) {
	parts := make([]string, len(a))
	params := make([]interface{}, 0, len(a))

	for i, c := range a {
		q, p, err := c.Build()
		if err != nil {
			return "", nil, err
		}
		parts[i] = q
		params = append(params, p...)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, " AND ")), params, nil
}

// Build builds a predicate. If the Pred's value implements the Builder
// interface, then the output of its Build method is used as the predicate's
// expression. Otherwise, the expression is set to a `?`.
func (c Pred) Build() (q string, p []interface{}, err error) {
	switch v := c.Val.(type) {
	case Builder:
		q, p, err = v.Build()
		if err != nil {
			return "", nil, err
		}

		q = fmt.Sprintf("%s%s(%s)", c.Col, c.Op, q)
	default:
		q = fmt.Sprintf("%s%s?", c.Col, c.Op)
		p = []interface{}{c.Val}
	}
	return q, p, nil
}

type predicates []Builder

func (w predicates) Build() (string, []interface{}, error) {
	var parts []string
	var params []interface{}

	for _, c := range w {
		part, param, err := c.Build()
		if err != nil {
			return "", nil, err
		}

		parts = append(parts, part)
		params = append(params, param...)

	}
	return strings.Join(parts, " AND "), params, nil
}
