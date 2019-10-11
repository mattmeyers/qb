package qb

import (
	"fmt"
	"strings"
)

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
	return fmt.Sprint("(", strings.Join(p, ", "), ")")
}
