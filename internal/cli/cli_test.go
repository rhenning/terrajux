package cli

import (
	"regexp"
	"testing"

	"github.com/rhenning/terrajux"
	"github.com/stretchr/testify/assert"
)

func TestCLI_ParseArgs(t *testing.T) {
	type fields struct {
		Args   []string
		Config *terrajux.Config
	}
	tests := []struct {
		name   string
		fields fields

		wantErr     bool
		wantErrType interface{}

		wantMessageRE string

		wantConfigValue *terrajux.Config
	}{
		{
			name: "noargs",
			fields: fields{
				Args:   []string{"terrajux"},
				Config: &terrajux.Config{Name: "terrajux"},
			},

			wantErr:     true,
			wantErrType: &ArgumentError{},
			wantMessageRE: regexp.QuoteMeta(
				`Usage: terrajux [options] <giturl> <v1ref> <v2ref> [subpath]`,
			),
		},
		{
			name: "versionflag",
			fields: fields{
				Args: []string{"terrajux", "-version"},
				Config: &terrajux.Config{
					Name:       "tjux",
					Version:    "0.0.2-test.1",
					ProjectURL: "irc://irc.efnet.org/greetz",
				},
			},

			wantErr:       true,
			wantErrType:   &VersionError{},
			wantMessageRE: regexp.QuoteMeta(`tjux 0.0.2-test.1 irc://irc.efnet.org/greetz`),
		},
		{
			name: "noflags+okargs",
			fields: fields{
				Args:   []string{"terrajux", "git@test.com:a/b.git", "v1.2.3", "main"},
				Config: &terrajux.Config{},
			},

			wantErr: false,
			wantConfigValue: &terrajux.Config{
				GitURL:   "git@test.com:a/b.git",
				GitRefV1: "v1.2.3",
				GitRefV2: "main",
			},
		},
		{
			name: "noflags+okargs+subpath",
			fields: fields{
				Args:   []string{"terrajux", "git@test.com:a/b.git", "v1.2.3", "main", "sub/dir"},
				Config: &terrajux.Config{},
			},

			wantErr: false,
			wantConfigValue: &terrajux.Config{
				GitURL:     "git@test.com:a/b.git",
				GitRefV1:   "v1.2.3",
				GitRefV2:   "main",
				GitSubpath: "sub/dir",
			},
		},
		{
			name: "clearcacheflag+okargs+subpath",
			fields: fields{
				Args:   []string{"terrajux", "-clearcache", "url", "v1", "v2", "sub/p"},
				Config: &terrajux.Config{},
			},

			wantErr: false,
			wantConfigValue: &terrajux.Config{
				CacheClear: true,
				GitURL:     "url",
				GitRefV1:   "v1",
				GitRefV2:   "v2",
				GitSubpath: "sub/p",
			},
		},
		{
			name: "difftoolflag+okargs",
			fields: fields{
				Args: []string{
					"terrajux",
					"-difftool",
					"opendiff {{.V1}} {{.V2}}",
					"foo", "bar", "baz",
				},
				Config: &terrajux.Config{},
			},

			wantConfigValue: &terrajux.Config{
				DiffTool:   "opendiff {{.V1}} {{.V2}}",
				CacheClear: false,
				GitURL:     "foo",
				GitRefV1:   "bar",
				GitRefV2:   "baz",
				GitSubpath: "",
			},
		},
		{
			name: "badflag",
			fields: fields{
				Args:   []string{"terrajux", "-bad=jawn", "url", "v1", "v2", "sub/p"},
				Config: terrajux.NewConfig(),
			},

			wantErr:       true,
			wantErrType:   &HelpError{},
			wantMessageRE: regexp.QuoteMeta(`Usage: terrajux`),
		},
		{
			name: "toomany",
			fields: fields{
				Args:   []string{"terrajux", "url", "v1", "v2", "sub/p", "nope"},
				Config: terrajux.NewConfig(),
			},

			wantErr:       true,
			wantErrType:   &ArgumentError{},
			wantMessageRE: regexp.QuoteMeta(`Usage: terrajux`),
		},
	}
	for _, tt := range tests {
		// capture range variable for t.Parallel()
		// forget this and you'll go mad chasing flaky test results
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			clii := New(tt.fields.Args, tt.fields.Config)

			gotMessage, err := clii.ParseArgs()

			if tt.wantErr {
				assert.Errorf(err, "CLI.ParseArgs() wantErr=%v, error=`%v`", tt.wantErr, err)
				if tt.wantErrType != nil {
					assert.IsType(
						tt.wantErrType, err,
						"CLI.ParseArgs wantErrType=%v, err=`%v`",
						tt.wantErrType, err,
					)
				}
			} else {
				assert.NoErrorf(err, "CLI.ParseArgs() wantErr=%v, err=`%v`", tt.wantErr, err)
			}

			if tt.wantMessageRE == "" {
				assert.Zerof(gotMessage, "CLI.ParseArgs() message=`%v`", gotMessage)
			} else {
				assert.Regexpf(
					regexp.MustCompile(tt.wantMessageRE), gotMessage,
					"CLI.ParseArgs() wantMessageRE=`%v`, message=`%v`",
					tt.wantMessageRE, gotMessage,
				)
			}

			if tt.wantConfigValue != nil {
				assert.Equalf(
					tt.wantConfigValue, tt.fields.Config,
					"CLI.ParseArgs() wantConfigValue=%+v, config=%+v",
					tt.wantConfigValue, clii.Config,
				)
			}
		})
	}
}
