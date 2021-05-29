package git

import (
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Clone(url string, ref string, dir string) error {
	opts := mkCloneOptions(url, plumbing.NewBranchReferenceName(ref))
	_, err := git.PlainClone(dir, false, opts)

	// is err of type git.NoMatchingRefSpecError?
	if _, ok := err.(git.NoMatchingRefSpecError); ok {
		opts = mkCloneOptions(url, plumbing.NewTagReferenceName(ref))
		_, err = git.PlainClone(dir, false, opts)
	}

	return err
}

func mkCloneOptions(url string, ref plumbing.ReferenceName) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           url,
		ReferenceName: ref,
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	}
}
