package terraform

import (
	"os"
	"path/filepath"
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

	assert.NoError(t, tfcli.Version(), "CLI.Version() should not error")
}

func TestCLIInit(t *testing.T) {
	t.Parallel()

	tdir := th.CreateTempDir(t)
	defer os.RemoveAll(tdir)

	err := th.WriteFile(t, filepath.Join(tdir, "test.tf"), testModule)
	assert.NoError(t, err, "write a TF module to tempdir prior to CLI.Init()")

	tfcli := NewCLI()
	err = tfcli.Init(tdir)

	assert.NoError(t, err, "CLI.Init() should not error")
}
