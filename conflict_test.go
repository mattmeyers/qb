package qb

import (
	"testing"
)

func TestOnConflict(t *testing.T) {
	// q, p, err := InsertInto("test_table").Col("a", 1).OnConflict(TargetColumn("a"), ActionDoNothing).String()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(q)
	// fmt.Println(p)

	// q, p, err = InsertInto("test_table").Col("a", 1).OnConflict(TargetColumn("a"), Update("").Set("a", 2)).String()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(q)
	// fmt.Println(p)

	// q, p, err = InsertInto("test_table").Col("a", 1).OnConflict(TargetConstraint("my_constraint"), Update("").Set("b", false).Where("c", ">", 5)).String()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(q)
	// fmt.Println(p)

	// wClause := whereClause{
	// 	clauses: []Cmp{
	// 		Cmp{col: "a", op: "=", val: "b"},
	// 		Cmp{col: "c", op: ">", val: 5},
	// 		Cmp{col: "d", op: "<>", val: false},
	// 	},
	// }
	// q, p, err = InsertInto("test_table").Col("a", 1).OnConflict(wClause, Update("").Set("b", false).Where("c", ">", 5)).String()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(q)
	// fmt.Println(p)
}
