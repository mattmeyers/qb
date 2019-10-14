package qb

import (
	"reflect"
	"testing"
)

// Because we're storing setPairs in a map, the order of the sets doesn't remain the same.
// If these tests fail, check manually. They're probably fine.
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
			want:    "UPDATE test_table SET a=?, c=?",
			want1:   []interface{}{"b", 1},
			wantErr: false,
		},
		{
			name:    "Update with where clause",
			query:   Update("test_table").Set("a", "b").Set("c", 1).Where("c", "=", "d").OrWhere("e", "<", 1).Where("f", "!=", false),
			want:    "UPDATE test_table SET a=?, c=? WHERE c=? OR e<? AND f!=?",
			want1:   []interface{}{"b", 1, "d", 1, false},
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
