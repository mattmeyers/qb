package qb

import (
	"reflect"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
)

func Test_selectQuery_String(t *testing.T) {
	tests := []struct {
		name    string
		query   *SelectQuery
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "Select *",
			query:   Select().From("test_table"),
			want:    "SELECT * FROM test_table",
			want1:   nil,
			wantErr: false,
		},
		{
			name: "Select with where clause",
			query: Select("a", "b").
				From("test_table").
				Where(
					Or{
						And{
							Pred{"c", "=", "d"},
							Pred{"m", ">", 5},
						},
						Pred{"e", "<", 1},
					}).
				Where(Pred{"f", "!=", false}),
			want:    "SELECT a, b FROM test_table WHERE ((c=? AND m>?) OR e<?) AND f!=?",
			want1:   []interface{}{"d", 5, 1, false},
			wantErr: false,
		},
		{
			name:    "Nested where query",
			query:   Select("id").From("test_table").Where(Pred{"id", " in ", Select("id").From("second_table")}),
			want:    "SELECT id FROM test_table WHERE id in (SELECT id FROM second_table)",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Having clause",
			query:   Select("id").From("test_table").GroupBy("id").Having(Pred{"COUNT(*)", ">", 2}),
			want:    "SELECT id FROM test_table GROUP BY id HAVING COUNT(*)>?",
			want1:   []interface{}{2},
			wantErr: false,
		},
		{
			name:    "Inner join",
			query:   Select().From("test_table").InnerJoin("second_table", S("test_table.id=second_table.test_table_id")),
			want:    "SELECT * FROM test_table INNER JOIN second_table ON test_table.id=second_table.test_table_id",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Limit and offset",
			query:   Select("a", "b").From("test_table").Limit(10).Offset(20),
			want:    "SELECT a, b FROM test_table LIMIT 10 OFFSET 20",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Group by",
			query:   Select("a", "b").From("test_table").GroupBy("a", "b"),
			want:    "SELECT a, b FROM test_table GROUP BY a, b",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Order by",
			query:   Select("a", "b").From("test_table").OrderBy("a", Asc).OrderBy("b", Desc),
			want:    "SELECT a, b FROM test_table ORDER BY a ASC, b DESC",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Missing table",
			query:   Select(),
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name:    "Select distinct",
			query:   Select("a", "b").Distinct().From("test_table"),
			want:    "SELECT DISTINCT a, b FROM test_table",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Select distinct on",
			query:   Select("a", "b", "c").Distinct("a").From("test_table").Distinct("b"),
			want:    "SELECT DISTINCT ON (a, b) a, b, c FROM test_table",
			want1:   nil,
			wantErr: false,
		},
		// {
		// 	name:    "Select from derived table",
		// 	query:   Select(""),
		// 	want:    "",
		// 	want1:   nil,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.query.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("insertQuery.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("insertQuery.String() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("insertQuery.String() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestThis(t *testing.T) {
	// t.Log(qbSelect())
	t.Log(squirrelSelect())
	// t.Log(dbrSelect())
}

func qbSelect() string {
	return Select().
		From("test_table").
		InnerJoin("second_table", S("test_table.id=second_table.test_table_id")).
		LeftJoin("third_table", S("second_table.third_id=third_table.id")).
		Limit(10).
		Offset(15).
		GroupBy("first_table.id").
		OrderBy("second_table.id", Asc).
		OrderBy("third_table.id", Desc).
		Where(Or{Pred{"a", "=", "b"}, Pred{"c", "=", "d"}, Pred{"e", "=", Select("id").From("sub_table").Where(Pred{"f", "=", 1})}}).
		String()
}

func squirrelSelect() string {
	q, _, _ := squirrel.Select("*").
		From("test_table").
		Join("second_table", "test_table.id=second_table.test_table_id").
		LeftJoin("third_table", "second_table.third_id=third_table.id").
		Limit(10).
		Offset(15).
		GroupBy("first_table.id").
		OrderBy("second_table.id", "third_table.id").
		Where(squirrel.Or{squirrel.Eq{"a": "b", "c": "d", "e": squirrel.Select("id").From("sub_table").Where(squirrel.Eq{"f": 1})}}).
		ToSql()
	return q
}

func dbrSelect() string {
	buf := dbr.NewBuffer()
	dbr.Select("*").
		From("test_table").
		Join("second_table", "test_table.id=second_table.test_table_id").
		LeftJoin("third_table", "second_table.third_id=third_table.id").
		Limit(10).
		Offset(15).
		GroupBy("first_table.id").
		OrderAsc("second_table.id").
		OrderDesc("third_table.id").
		Where(dbr.Or(dbr.Eq("a", "b"), dbr.Eq("c", "d"), dbr.Eq("e", "f"))).
		Build(dialect.PostgreSQL, buf)
	return buf.String()

}

func Benchmark_SelectQuery(b *testing.B) {
	tests := []struct {
		name string
		fun  func() string
	}{
		{"qb", qbSelect},
		{"sqirrel", squirrelSelect},
		{"dbr", dbrSelect},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				test.fun()
			}
		})
	}

	for i := 0; i < b.N; i++ {
	}
}
