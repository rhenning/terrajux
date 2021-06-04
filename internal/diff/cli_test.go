package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{})

	assert.NoError(err, "NewRunner() should not error")
	assert.Equal("/bin/sh", dr.Options.Shell, "NewRunner() should use /bin/sh by default")
}

func TestNewRunnerWithOptions(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{
		Dir:   "/tmp/foo",
		Shell: "/usr/local/bin/zsh",
	})

	assert.NoError(err, "NewRunner() should not error with options")
	assert.Equal("/usr/local/bin/zsh", dr.Options.Shell, "NewRunner() should use Shell override")
	assert.Equal("/tmp/foo", dr.Options.Dir, "NewRunner() should use Dir override")
}

func TestDiffRun(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{})

	assert.NoError(err, "NewRunner() should not error")

	err = dr.Run("testdata/a", "testdata/b")

	assert.NoError(err, "Runner.Run(v1, v2) should not error")
}

func TestDiffRunInDir(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{Dir: "testdata"})

	assert.NoError(err, "NewRunner() should not error")

	err = dr.Run("a", "b")

	assert.NoError(err, "Runner.Run(v1, v2) should not error")
}

func TestDiffRunError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{})

	assert.NoError(err, "NewRunner() should not error")

	err = dr.Run("testdata/this-is-a-farce", "testdata/b")

	assert.Error(err, "Runner.Run(v1, v2) should error with bad diff args")
}

func TestFormatCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{
		Shell:           "/bin/foo",
		CommandTemplate: "zork -v2={{ .V2 }} -v1 {{ .V1 }}",
	})

	assert.NoErrorf(err, "NewRunner() should not error")

	err = dr.formatCommand("foo", "bar")

	assert.NoError(err, "Runner.formatCommand() should not error")
	assert.Equal("zork -v2=bar -v1 foo", dr.command, "Runner.formatCommand() should correctly format")
}

func TestFormatCommandError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dr, err := NewRunner(&RunnerOptions{
		Shell:           "/bin/foo",
		CommandTemplate: "test {{ .BadTemplate }}",
	})

	assert.NoErrorf(err, "NewRunner() should not error")
	assert.Error(dr.formatCommand("a", "b"), "Runner.formatCommand() should error with bad template")
}
