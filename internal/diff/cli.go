package diff

import (
	"bytes"
	"fmt"
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
		"{{.V1}}", // template replaced by path for ref v1
		"{{.V2}}", // template replaced by path for ref v2
	}
)

type ToolOptions struct {
	Dir             string
	Env             []string
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

	ct, err := template.New("diffcmd").Parse(
		quoteTemplate(tool.Options.CommandTemplate),
	)

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

	// args are composed from values defined in this package or from config
	// directly supplied by the user/administrator. still, we quote the
	// parameterized input to guard against shell expansion.
	cmd := exec.Command(tool.Options.Shell, "-c", tool.command) // #nosec G204

	cmd.Dir = tool.Options.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// difftool might need a few env vars to work properly
	for _, ekey := range []string{"HOME", "PATH", "TMPDIR"} {
		if eval, ok := os.LookupEnv(ekey); ok {
			tool.Options.Env = append(tool.Options.Env, fmt.Sprintf("%s=%s", ekey, eval))
		}
	}

	cmd.Env = tool.Options.Env

	err := cmd.Run()

	if exErr, ok := err.(*exec.ExitError); ok {
		if exErr.ExitCode() == 0 || exErr.ExitCode() == 1 {
			return nil
		}
	}

	return err
}

func (c *Tool) formatCommand(v1path string, v2path string) error {
	var command bytes.Buffer

	ctargs := commandTemplateArgs{v1path, v2path}

	err := c.ct.Execute(&command, ctargs)

	c.command = command.String()

	return err
}

func quoteTemplate(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(s, "{{", "'{{"),
		"}}", "}}'",
	)
}
