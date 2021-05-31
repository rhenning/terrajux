package git

import (
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	cloneBare         = false
	cloneDepth        = 1
	cloneSingleBranch = true
)

type Cloner interface {
	Clone() error
}

type Git struct{}

func New() *Git {
	return &Git{}
}

func (g *Git) Clone(url string, ref string, dir string) error {
	opts := g.mkCloneOptions(url, plumbing.NewBranchReferenceName(ref))
	_, err := git.PlainClone(dir, cloneBare, opts)

	// is err of type git.NoMatchingRefSpecError?
	if _, ok := err.(git.NoMatchingRefSpecError); ok {
		opts = g.mkCloneOptions(url, plumbing.NewTagReferenceName(ref))
		_, err = git.PlainClone(dir, cloneBare, opts)
	}

	return err
}

func (g *Git) mkCloneOptions(url string, ref plumbing.ReferenceName) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           url,
		ReferenceName: ref,
		SingleBranch:  cloneSingleBranch,
		Depth:         cloneDepth,
		//Progress:      os.Stdout, // for more info
	}
}
