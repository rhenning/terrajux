package cache

import (
	"os"
	"path/filepath"
	"testing"

	th "github.com/rhenning/terrajux/internal/testhelp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	t.Parallel()

	c := New("/tmp/a/z")

	assert.Equal(t, "/tmp/a/z", c.Dir, "New(dir) should set Cache.Dir")
}

func TestEnsureCache(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	tdir := th.CreateTempDir(t)
	defer os.RemoveAll(tdir)

	c := New(filepath.Join(tdir, "cache"))
	assert.NoError(c.Ensure(), "Cache.Ensure() should not error for new Cache.Dir")
	defer os.RemoveAll(c.Dir)

	assert.DirExists(c.Dir, "Cache.Dir should exist and be a directory")
	assert.NoError(c.Ensure(), "Cache.Ensure() should not error when Cache.Dir exists")

	err := th.WriteFile(t, filepath.Join(c.Dir, "junk.dat"), "x")
	assert.NoError(err, "writing a new file in Cache.Dir should not error")
}

func TestClearCache(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	tdir := th.CreateTempDir(t)
	defer os.RemoveAll(tdir)

	c := New(filepath.Join(tdir, "cache"))
	require.NoError(c.Ensure(), "ensure Cache.Dir exists prior to Cache.Clear()")
	defer os.RemoveAll(c.Dir)

	err := th.WriteFile(t, filepath.Join(c.Dir, "junk.dat"), "x")
	assert.NoError(err, "write a new file to Cache.Dir prior to Cache.Clear()")

	assert.DirExists(c.Dir, "Cache.Dir should exist and be a directory")
	assert.NoError(c.Clear(), "Cache.Clear() should not error")

	ok, err := th.DirIsEmpty(c.Dir)
	assert.NoErrorf(err, "Cache.Dir empty check should not error")
	assert.True(ok, "Cache.Dir should be empty after cache.Clear()")

	err = th.WriteFile(t, filepath.Join(c.Dir, "junk.dat"), "x")
	assert.NoError(err, "writing a new file in Cache.Dir should not error after Cache.Clear()")
}
