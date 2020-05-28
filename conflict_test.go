package qb

import (
	"fmt"
	"testing"
)

func TestOnConflict(t *testing.T) {
	q, p, err := InsertInto("test_table").Col("a", 1).OnConflict(TargetColumn("a"), ActionDoNothing).Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(q)
	fmt.Println(p)

	q, p, err = InsertInto("test_table").Col("a", 1).OnConflict(TargetColumn("a"), Update("").Set("a", 2)).Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(q)
	fmt.Println(p)

	q, p, err = InsertInto("test_table").Col("a", 1).OnConflict(TargetConstraint("my_constraint"), Update("").Set("b", false).Where(Gt("c", 5))).Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(q)
	fmt.Println(p)
}
