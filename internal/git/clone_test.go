package git

import (
	"errors"
	"os"
	"testing"

	th "github.com/rhenning/terrajux/internal/testhelp"
	"github.com/stretchr/testify/assert"
)

func TestGit_Clone(t *testing.T) {
	type args struct {
		url string
		ref string
		dir string
	}
	tests := []struct {
		name        string
		expectedErr error
		args        args
	}{
		{
			name:        "branch",
			expectedErr: nil,
			args: args{
				url: "https://github.com/terraform-aws-modules/terraform-aws-iam.git",
				ref: "master",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "fully qualified branch reference",
			expectedErr: nil,
			args: args{
				url: "https://github.com/terraform-aws-modules/terraform-aws-iam.git",
				ref: "refs/heads/master",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "nonexistent fully qualified branch reference",
			expectedErr: errors.New("couldn't find remote ref \"refs/heads/foo\""),
			args: args{
				url: "https://github.com/terraform-aws-modules/terraform-aws-iam.git",
				ref: "refs/heads/foo",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "tag",
			expectedErr: nil,
			args: args{
				url: "https://github.com/go-git/go-git.git",
				ref: "v4.0.0",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "fully qualified tag reference",
			expectedErr: nil,
			args: args{
				url: "https://github.com/go-git/go-git.git",
				ref: "refs/tags/v4.0.0",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "nonexistent fully qualified tag reference",
			expectedErr: errors.New("couldn't find remote ref \"refs/tags/v2000.0.0\""),
			args: args{
				url: "https://github.com/go-git/go-git.git",
				ref: "refs/tags/v2000.0.0",
				dir: th.CreateTempDir(t),
			},
		},
		{
			name:        "nonexistent reference containg '/'",
			expectedErr: errors.New("couldn't find remote ref \"refs/tags/foo/master\""),
			args: args{
				url: "https://github.com/terraform-aws-modules/terraform-aws-iam.git",
				ref: "foo/master",
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

			git := New()
			err := git.Clone(tt.args.url, tt.args.ref, tt.args.dir)

			switch tt.expectedErr {
			case nil:
				assert.NoErrorf(t, err, "Clone(%+v) error: %v", tt.args, err)
			default:
				assert.Equal(t, err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestURLPath(t *testing.T) {
	type args struct {
		url string
		ref string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{
			name: "standard+scheme+noport+nouser",
			args: args{
				url: "https://test.com/foo/bar.git",
				ref: "main",
			},

			wantPath: "test.com/foo/bar.git/main",
		},
		{
			name: "ssh+noscheme+user",
			args: args{
				url: "git@test.com:foo/bar.git",
				ref: "v2.0.0",
			},

			wantPath: "test.com/foo/bar.git/v2.0.0",
		},
	}
	for _, tt := range tests {
		// capture range variable for t.Parallel()
		// forget this and you'll go mad chasing flaky test results
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			gotPath := URLPath(tt.args.url, tt.args.ref)

			if tt.wantPath != "" {
				assert.Equalf(
					tt.wantPath, gotPath,
					"URLPath() wantPath=`%v`, gotPath=`%v`",
					tt.wantPath, gotPath,
				)
			}
		})
	}
}
