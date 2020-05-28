package qb

import (
	"reflect"
	"testing"
)

func Test_joinClause_String(t *testing.T) {
	tests := []struct {
		name    string
		jc      joins
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name:    "Single inner join",
			jc:      joins{join{innerJoin, "b", "a.id=b.a_id"}},
			want:    "INNER JOIN b ON a.id=b.a_id",
			want1:   nil,
			wantErr: false,
		},
		{
			name: "Multiple inner joins",
			jc: joins{
				join{innerJoin, "b", "a.id=b.a_id"},
				join{innerJoin, "c", "b.id=c.b_id"},
			},
			want:    "INNER JOIN b ON a.id=b.a_id INNER JOIN c ON b.id=c.b_id",
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.jc.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("joinClause.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("joinClause.String() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("joinClause.String() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
