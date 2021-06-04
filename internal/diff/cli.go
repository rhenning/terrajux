package diff

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	defaultShell = "/bin/sh"

	// elements will be joined with spaces before expansion & execution by defaultShell
	defaultCommandTemplate = []string{
		"git",
		"-c", "advice.detachedHead=", // apparently git interprets foo= as false
		"diff",
		"--color=auto",
		"--no-ext-diff",
		"--no-index",
		"--ignore-all-space",
		"{{ .V1 }}", // template replaced by path for ref v1
		"{{ .V2 }}", // template replaced by path for ref v2
	}
)

type RunnerOptions struct {
	Dir             string
	Shell           string
	CommandTemplate string
}

type Runner struct {
	Options *RunnerOptions
	ct      *template.Template
	command string
}

type commandTemplateArgs struct {
	V1 string
	V2 string
}

func NewRunner(opts *RunnerOptions) (*Runner, error) {
	runner := &Runner{
		Options: &RunnerOptions{
			Dir:             opts.Dir,
			Shell:           defaultShell,
			CommandTemplate: strings.Join(defaultCommandTemplate, " "),
		},
	}

	if opts.Shell != "" {
		runner.Options.Shell = opts.Shell
	}

	if opts.CommandTemplate != "" {
		runner.Options.CommandTemplate = opts.CommandTemplate
	}

	ct, err := template.New("diffcmd").Parse(runner.Options.CommandTemplate)
	runner.ct = ct

	return runner, err
}

func (c *Runner) Run(v1path string, v2path string) error {
	for _, p := range []string{v1path, v2path} {
		dirp := filepath.Join(c.Options.Dir, p)

		if _, err := os.Stat(dirp); err != nil {
			return err
		}
	}

	if err := c.formatCommand(v1path, v2path); err != nil {
		return err
	}

	cmd := exec.Command(c.Options.Shell, "-c", c.command)

	cmd.Dir = c.Options.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if exErr, ok := err.(*exec.ExitError); ok {
		if exErr.ExitCode() == 0 || exErr.ExitCode() == 1 {
			return nil
		}
	}

	return err
}

func (c *Runner) formatCommand(v1path string, v2path string) error {
	ctargs := commandTemplateArgs{v1path, v2path}
	var command bytes.Buffer

	err := c.ct.Execute(&command, ctargs)
	c.command = command.String()

	return err
}
