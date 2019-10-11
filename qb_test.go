package qb

import (
	"testing"
)

func TestGeneratePlaceholders(t *testing.T) {
	type args struct {
		symbol string
		num    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Negative repeat",
			args: args{"?", -1},
			want: "()",
		},
		{
			name: "Repeat zero times",
			args: args{"?", 0},
			want: "()",
		},
		{
			name: "Repeat once",
			args: args{"?", 1},
			want: "(?)",
		},
		{
			name: "Repeat twice",
			args: args{"?", 2},
			want: "(?, ?)",
		},
		{
			name: "Repeat thrice",
			args: args{"?", 3},
			want: "(?, ?, ?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePlaceholders(tt.args.symbol, tt.args.num); got != tt.want {
				t.Errorf("GeneratePlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}
