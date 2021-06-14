package git

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	cloneBare         = false
	cloneDepth        = 1
	cloneSingleBranch = true
)

type Cloner interface {
	Clone(url string, ref string, dir string) error
}

type Git struct{}

func New() *Git {
	return &Git{}
}

func (g *Git) Clone(url string, ref string, dir string) error {
	// is the ref a fully qualified ref such as 'refs/tags/<ref>'?
	if strings.HasPrefix(ref, "refs/") {
		opts := g.mkCloneOptions(url, plumbing.ReferenceName(ref))
		_, err := git.PlainClone(dir, cloneBare, opts)
		if err != nil {
			return err
		}

		return nil
	}

	// attempt to to clone a branch reference such as 'refs/head/<ref>'
	opts := g.mkCloneOptions(url, plumbing.NewBranchReferenceName(ref))
	_, err := git.PlainClone(dir, cloneBare, opts)

	// is err of type git.NoMatchingRefSpecError?
	// attempt to to clone a tag reference such as 'refs/tags/<ref>'
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

func URLPath(u string, ref string) (path string) {
	up, err := url.Parse(u)

	if err != nil {
		return urlPathSSH(u, ref)
	}

	path = fmt.Sprintf("%s/%s/%s/%s/%s", up.Scheme, up.Hostname(), up.Port(), up.Path, ref)
	return filepath.Clean(path)
}

func urlPathSSH(u string, ref string) (path string) {
	ss := strings.SplitN(u, "@", 2)
	hostpath := ss[len(ss)-1]
	hostpath = strings.Replace(hostpath, ":", "/", 1)
	path = fmt.Sprintf("%s/%s", hostpath, ref)
	return filepath.Clean(path)
}
