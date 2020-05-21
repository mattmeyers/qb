package qb

import (
	"fmt"
	"strings"
)

type conflictResolver struct {
	target interface{}
	action interface{}
}

type TargetColumn string
type TargetConstraint string

type actionDoNothing string

const ActionDoNothing = actionDoNothing("NOTHING")

func (c *conflictResolver) String() (string, []interface{}, error) {
	var sb strings.Builder
	var params []interface{}
	sb.WriteString("ON CONFLICT ")

	switch v := c.target.(type) {
	case string:
		fmt.Fprintf(&sb, "(%s)", v)
	case TargetColumn:
		fmt.Fprintf(&sb, "(%s)", string(v))
	case TargetConstraint:
		fmt.Fprintf(&sb, "ON CONSTRAINT %s", string(v))
	case whereClause:
		q, p, _ := v.SQL()
		fmt.Fprintf(&sb, "WHERE %s", q)
		params = append(params, p...)
	default:
		return "", nil, ErrInvalidConflictTarget
	}

	sb.WriteString(" DO ")

	switch v := c.action.(type) {
	case actionDoNothing:
		sb.WriteString(string(v))
	case *updateQuery:
		q, p, err := v.string(false)
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(q)
		params = append(params, p...)
	default:
		return "", nil, ErrInvalidConflictAction
	}

	return sb.String(), params, nil
}
