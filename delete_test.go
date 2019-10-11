package qb

import (
	"reflect"
	"testing"
)

func Test_deleteQuery_String(t *testing.T) {
	tests := []struct {
		name    string
		query   *deleteQuery
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "Delete all",
			query:   DeleteFrom("test_table"),
			want:    "DELETE FROM test_table",
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "Delete with where clause",
			query:   DeleteFrom("test_table").Where("c", "=", "d").OrWhere("e", "<", 1).Where("f", "!=", false),
			want:    "DELETE FROM test_table WHERE c=? OR e<? AND f!=?",
			want1:   []interface{}{"d", 1, false},
			wantErr: false,
		},
		{
			name:    "Missing table",
			query:   DeleteFrom(""),
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
