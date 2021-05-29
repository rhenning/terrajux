package git

import (
	"os"
	"testing"

	th "github.com/rhenning/terrajux/internal/testhelp"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	type args struct {
		url string
		ref string
		dir string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "branch",
			args: args{
				url: "https://github.com/terraform-aws-modules/terraform-aws-iam.git",
				ref: "master",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name: "tag",
			args: args{
				url: "https://github.com/go-git/go-git.git",
				ref: "v4.0.0",
				dir: th.CreateTempDir(t),
			},
		},
	}
	for _, tt := range tests {
		// capture range variable for t.Parallel()
		// forget this and you'll go mad chasing flaky test results
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer os.RemoveAll(tt.args.dir)

			err := Clone(tt.args.url, tt.args.ref, tt.args.dir)
			assert.NoErrorf(t, err, "Clone(%+v) error: %v", tt.args, err)
		})
	}
}
