package qb

import (
	"reflect"
	"testing"
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
