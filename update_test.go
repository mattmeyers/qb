package qb

import (
	"reflect"
	"testing"
)

func Test_updateQuery_String(t *testing.T) {
	tests := []struct {
		name    string
		query   *updateQuery
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "Simple update",
			query:   Update("test_table").Set("a", "b").Set("c", 1),
			want:    `UPDATE "test_table" SET a=?, c=?`,
			want1:   []interface{}{"b", 1},
			wantErr: false,
		},
		{
			name:    "Update with EXCLUDED value",
			query:   Update("test_table").Set("a", "b").Set("c", Excluded("c")),
			want:    `UPDATE "test_table" SET a=?, c=EXCLUDED.c`,
			want1:   []interface{}{"b"},
			wantErr: false,
		},
		{
			name:    "Update with where clause",
			query:   Update("test_table").Set("a", "b").Set("c", 1).Where(Eq("c", "d")).Where(Neq("f", false)),
			want:    `UPDATE "test_table" SET a=?, c=? WHERE c=? AND f!=?`,
			want1:   []interface{}{"b", 1, "d", false},
			wantErr: false,
		},
		{
			name:    "Missing table",
			query:   Update("").Set("a", "b"),
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name:    "Missing set pair",
			query:   Update("test_table"),
			want:    "",
			want1:   nil,
			wantErr: true,
		},
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
