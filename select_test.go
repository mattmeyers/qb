package qb

import (
	"reflect"
	"testing"

	"github.com/masterminds/squirrel"
)

func Test_selectQuery_String(t *testing.T) {
	tests := []struct {
		name    string
		query   *selectQuery
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
			name:    "Select with where clause",
			query:   Select("a", "b").From("test_table").Where("c", "=", "d").OrWhere("e", "<", 1).Where("f", "!=", false),
			want:    "SELECT a, b FROM test_table WHERE c=? OR e<? AND f!=?",
			want1:   []interface{}{"d", 1, false},
			wantErr: false,
		},
		{
			name:    "Inner join",
			query:   Select().From("test_table").InnerJoin("second_table", "test_table.id=second_table.test_table_id"),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.query.String()
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

func qbSelect() {
	Select().
		From("test_table").
		InnerJoin("second_table", "test_table.id=second_table.test_table_id").
		LeftJoin("third_table", "second_table.third_id=third_table.id").
		Limit(10).
		Offset(15).
		GroupBy("first_table.id").
		OrderBy("second_table.id", Asc).
		OrderBy("third_table.id", Desc).
		String()
}

func squirrelSelect() {
	squirrel.Select("*").
		From("test_table").
		Join("second_table", "test_table.id=second_table.test_table_id").
		LeftJoin("third_table", "second_table.third_id=third_table.id").
		Limit(10).
		Offset(15).
		GroupBy("first_table.id").
		OrderBy("second_table.id", "third_table.id").ToSql()
}

func Benchmark_SelectQuery(b *testing.B) {
	tests := []struct {
		name string
		fun  func()
	}{
		{"qb", qbSelect},
		{"sqirrel", squirrelSelect},
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
