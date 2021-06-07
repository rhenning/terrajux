package testhelp

import (
	"regexp"
	"testing"
)

func TestContainsRegexp(t *testing.T) {
	type args struct {
		c []string
		r *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match",
			args: args{
				c: []string{"THIS=that", "NEAR=far"},
				r: regexp.MustCompile("^NEAR="),
			},
			want: true,
		},
		{
			name: "nomatch",
			args: args{
				c: []string{"BIGBIRD=yellow", "OSCAR=green"},
				r: regexp.MustCompile("^GROVER="),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsRegexp(tt.args.c, tt.args.r)
			if got != tt.want {
				t.Errorf(
					"containsRegexp(c, r) want = %v, got = %v, c = %+v, r = %+v",
					tt.want, got, tt.args.c, tt.args.r,
				)
			}
		})
	}
}
