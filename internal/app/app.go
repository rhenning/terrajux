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
	Cache     cache.Initializer
	Config    terrajux.Config
}

func New(gc git.Cloner, ti terraform.Initer, dr diff.Runner, ci cache.Initializer, tc terrajux.Config) *App {
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
		dir := filepath.Join(a.Config.CacheDir, git.URLPath(a.Config.GitURL, v))
		clonedirs[i] = dir

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Cloning %v @%v\n", a.Config.GitURL, v)
			if err = a.Git.Clone(a.Config.GitURL, v, dir); err != nil {
				return err
			}
		} else {
			fmt.Printf("Found %v @%v in cache. Skipping clone.\n", a.Config.GitURL, v)
		}
	}

	for _, v := range clonedirs {
		fmt.Printf("Running `terraform init -backend=false` in %v`\n", filepath.Join(v, a.Config.GitSubpath))
		if err = a.Terraform.Init(filepath.Join(v, a.Config.GitSubpath)); err != nil {
			return err
		}
	}

	err = a.Diff.Run(
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

	dr, err := diff.NewRunner(&diff.RunnerOptions{
		CommandTemplate: config.DiffTool,
	})

	app.Diff = *dr

	return app, err
}
