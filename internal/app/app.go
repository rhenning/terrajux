package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rhenning/terrajux"
	"github.com/rhenning/terrajux/internal/cache"
	"github.com/rhenning/terrajux/internal/diff"
	"github.com/rhenning/terrajux/internal/git"
	"github.com/rhenning/terrajux/internal/terraform"
)

type App struct {
	Git       git.Cloner
	Terraform terraform.Initer
	Diff      diff.Runner
	Cache     cache.Manager
	Config    terrajux.Config
}

func New(gc git.Cloner, ti terraform.Initer, dr diff.Runner, ci cache.Manager, tc terrajux.Config) *App {
	return &App{
		Git:       gc,
		Terraform: ti,
		Diff:      dr,
		Cache:     ci,
		Config:    tc,
	}
}

func (a *App) Run() (err error) {
	var clonedirs [2]string

	if err = a.Cache.Ensure(); err != nil {
		return err
	}

	if a.Config.CacheClear {
		fmt.Printf("Clearing cache dir %v\n", a.Config.CacheDir)
		if err = a.Cache.Clear(); err != nil {
			return err
		}
	}

	for i, v := range []string{a.Config.GitRefV1, a.Config.GitRefV2} {
		dir := git.URLPath(a.Config.GitURL, v)
		clonedirs[i] = dir

		if a.Cache.HasKey(dir) {
			fmt.Printf("Found %q in cache. Skipping clone.\n", dir)
		} else {
			fmt.Printf("Cloning %v @%v\n", a.Config.GitURL, v)
			abspath, _ := a.Cache.GetAbsKeyPath(dir)
			fmt.Print(abspath)

			if err = a.Git.Clone(a.Config.GitURL, v, abspath); err != nil {
				return err
			}
		}
	}

	for _, v := range clonedirs {
		abspath, _ := a.Cache.GetAbsKeyPath(v)
		abspath = filepath.Join(abspath, a.Config.GitSubpath)

		fmt.Printf("Running `terraform init -backend=false` in %v`\n", abspath)
		if err = a.Terraform.Init(abspath); err != nil {
			return err
		}
	}

	if err = os.Chdir(a.Config.CacheDir); err != nil {
		return err
	}

	err = a.Diff.Run(
		// using a relative path here cleans up unified diff output a bit
		filepath.Join(clonedirs[0], a.Config.GitSubpath),
		filepath.Join(clonedirs[1], a.Config.GitSubpath),
	)

	return err
}

func NewDefaultWiring(config *terrajux.Config) (app *App, err error) {
	app = &App{}

	app.Config = *config
	app.Git = git.New()
	app.Terraform = terraform.NewCLI()
	app.Cache = cache.New(config.CacheDir)

	difftool, err := diff.NewTool(&diff.ToolOptions{
		CommandTemplate: config.DiffTool,
	})

	app.Diff = difftool

	return app, err
}
