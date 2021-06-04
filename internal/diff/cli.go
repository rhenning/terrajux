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

type ToolOptions struct {
	Dir             string
	Shell           string
	CommandTemplate string
}

type Tool struct {
	Options *ToolOptions
	ct      *template.Template
	command string
}

type Runner interface {
	Run(v1path string, v2path string) error
}

type commandTemplateArgs struct {
	V1 string
	V2 string
}

func NewTool(opts *ToolOptions) (*Tool, error) {
	tool := &Tool{
		Options: &ToolOptions{
			Dir:             opts.Dir,
			Shell:           defaultShell,
			CommandTemplate: strings.Join(defaultCommandTemplate, " "),
		},
	}

	if opts.Shell != "" {
		tool.Options.Shell = opts.Shell
	}

	if opts.CommandTemplate != "" {
		tool.Options.CommandTemplate = opts.CommandTemplate
	}

	ct, err := template.New("diffcmd").Parse(tool.Options.CommandTemplate)
	tool.ct = ct

	return tool, err
}

func (tool *Tool) Run(v1path string, v2path string) error {
	for _, p := range []string{v1path, v2path} {
		dirp := filepath.Join(tool.Options.Dir, p)

		if _, err := os.Stat(dirp); err != nil {
			return err
		}
	}

	if err := tool.formatCommand(v1path, v2path); err != nil {
		return err
	}

	cmd := exec.Command(tool.Options.Shell, "-c", tool.command)

	cmd.Dir = tool.Options.Dir
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

func (c *Tool) formatCommand(v1path string, v2path string) error {
	ctargs := commandTemplateArgs{v1path, v2path}
	var command bytes.Buffer

	err := c.ct.Execute(&command, ctargs)
	c.command = command.String()

	return err
}
