package qb

import (
	"reflect"
	"testing"
)

func Test_insertQuery_String(t *testing.T) {
	tests := []struct {
		name    string
		query   *insertQuery
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "Basic insert",
			query:   InsertInto("test_table").Col("a", "c").Col("b", "d"),
			want:    "INSERT INTO test_table (a, b) VALUES (?, ?)",
			want1:   []interface{}{"c", "d"},
			wantErr: false,
		},
		{
			name:    "Multiple column insert",
			query:   InsertInto("test_table").Cols([]string{"a", "b"}, []interface{}{"c", "d"}...),
			want:    "INSERT INTO test_table (a, b) VALUES (?, ?)",
			want1:   []interface{}{"c", "d"},
			wantErr: false,
		},
		{
			name:    "Staggered insert",
			query:   InsertInto("test_table").Cols([]string{"a", "b"}, []interface{}{"c", "d"}...).Col("e", 1),
			want:    "INSERT INTO test_table (a, b, e) VALUES (?, ?, ?)",
			want1:   []interface{}{"c", "d", 1},
			wantErr: false,
		},
		{
			name:    "Missing table",
			query:   InsertInto("").Col("a", "b"),
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name:    "Col-Val mismatch",
			query:   InsertInto("test_table").Cols([]string{}, "b"),
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
