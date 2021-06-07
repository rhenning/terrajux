package terraform

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	th "github.com/rhenning/terrajux/internal/testhelp"
	"github.com/stretchr/testify/assert"
)

const testModule = `
module "consul" {
  source = "github.com/terraform-aws-modules/terraform-aws-iam//modules/iam-user?ref=v4.1.0"
  name   = "alice"
}
`

func TestCLIVersion(t *testing.T) {
	t.Parallel()

	tfcli := NewCLI()
	err := tfcli.Version()

	assert.NoErrorf(t, err, "CLI.Version() should not error: %+v", err)
}

func TestCLIInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tdir := th.CreateTempDir(t)
	defer os.RemoveAll(tdir)

	err := th.WriteFile(t, filepath.Join(tdir, "test.tf"), testModule)
	assert.NoError(err, "write a TF module to tempdir prior to CLI.Init()")

	testevk := fmt.Sprintf("SEKRET%s", strconv.Itoa(rand.Int()))

	if err = os.Setenv(testevk, "flibbertygibbets"); err != nil {
		t.Logf("Couldn't set test environment var %s: %v", testevk, err)
	}

	defer func() {
		if derr := os.Unsetenv(testevk); derr != nil {
			err = derr
			fmt.Printf("%+v\n", err)
		}
	}()

	tfcli := NewCLI()
	err = tfcli.Init(tdir)

	assert.NoErrorf(err, "CLI.Init() should not error: %+v", err)
	assert.False(
		th.ContainsRegexp(tfcli.Env, regexp.MustCompile(fmt.Sprintf("^%s=", testevk))),
		"Expected %+v not to contain SEKRET* env var",
		tfcli.Env,
	)
	assert.True(
		th.ContainsRegexp(tfcli.Env, regexp.MustCompile("^HOME=")),
		"Expected %+v to contain HOME env var",
		tfcli.Env,
	)
}
