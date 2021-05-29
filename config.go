package terrajux

import (
	"fmt"
	"os/user"
	"path/filepath"
)

const Name = "terrajux"

var Version = "dev"

const Usage = `
https://github.com/rhenning/terrajux

terrajux diffs the source code of a terraform project stored in a git repo-
sitory, along with the source of all of its transitive module dependencies.

usage:

	terrajux [-clean] giturl ref1 ref2 [subpath]
	terrajux -clean

		giturl:  a git-compatible url
		ref1:    a starting git reference (tag, branch, or sha1)
		ref2:    an ending git reference
		subpath: an optional subpath of the repository containing the
					terraform module to initialize and compare

	   -clean:   wipe terrajux's git checkout and module cache 
`

func ConfigDir() string {
	u, err := user.Current()

	if err != nil {
		panic("couldn't lookup current user x.x")
	}

	return filepath.Join(u.HomeDir, fmt.Sprintf(".%s", Name))
}

func CacheDir() string {
	return filepath.Join(ConfigDir(), "cache")
}
