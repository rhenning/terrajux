package diff

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"testing"

	th "github.com/rhenning/terrajux/internal/testhelp"
	"github.com/stretchr/testify/assert"
)

func TestNewTool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{})

	assert.NoError(err, "NewTool() should not error")
	assert.Equal("/bin/sh", tool.Options.Shell, "NewTool() should use /bin/sh by default")
}

func TestNewToolWithOptions(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{
		Dir:   "/tmp/foo",
		Shell: "/usr/local/bin/zsh",
	})

	assert.NoError(err, "NewTool() should not error with options")
	assert.Equal("/usr/local/bin/zsh", tool.Options.Shell, "NewTool() should use Shell override")
	assert.Equal("/tmp/foo", tool.Options.Dir, "NewTool() should use Dir override")
}

func TestDiffRun(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{})

	assert.NoError(err, "NewTool() should not error")

	err = tool.Run("testdata/a", "testdata/b")

	assert.NoError(err, "Tool.Run(v1, v2) should not error")
}

func TestDiffRunInDir(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{Dir: "testdata"})

	assert.NoError(err, "NewTool() should not error")

	err = tool.Run("a", "b")

	assert.NoError(err, "Tool.Run(v1, v2) should not error")
}

func TestDiffRunError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{})

	assert.NoError(err, "NewTool() should not error")

	err = tool.Run("testdata/this-is-a-farce", "testdata/b")

	assert.Error(err, "Tool.Run(v1, v2) should error with bad diff args")
}

func TestToolCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{
		// nonportable: eventually inject a mock command runner
		Shell:           "/bin/bash",
		CommandTemplate: "HISTFILE=/dev/null echo -v2={{.V2}} -v1 {{.V1}}",
	})

	assert.NoErrorf(err, "NewTool() should not error")

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

	err = tool.formatCommand("foo", "bar")

	assert.NoError(err, "Tool.formatCommand() should not error")
	assert.Equal(
		"HISTFILE=/dev/null echo -v2='bar' -v1 'foo'",
		tool.command,
		"Tool.formatCommand() should correctly format and quote params",
	)

	err = tool.Run("testdata/a", "testdata/b")

	assert.NoError(err, "Tool.Run(v1, v2) should not error")
	assert.Falsef(
		th.ContainsRegexp(
			tool.Options.Env,
			regexp.MustCompile(fmt.Sprintf("^%s=", testevk)),
		),
		"",
	)
}

func TestFormatCommandError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tool, err := NewTool(&ToolOptions{
		Shell:           "/bin/foo",
		CommandTemplate: "test {{.BadTemplate}}",
	})

	assert.NoErrorf(err, "NewTool() should not error")
	assert.Error(tool.formatCommand("a", "b"), "Tool.formatCommand() should error with bad template")
}
