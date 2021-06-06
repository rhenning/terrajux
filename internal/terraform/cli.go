package terraform

import (
	"os"
	"os/exec"
)

const (
	ExecCommand = "terraform"

	tfInitCommand    = "init"
	tfVersionCommand = "version"
)

var (
	globalEnv   = []string{"TF_IN_AUTOMATION=1"}
	globalFlags = []string{}

	tfInitCommandFlags = []string{
		"-input=false",
		"-backend=false",
	}

	tfVersionCommandFlags = []string{}
)

type Initer interface {
	Init(dir string) error
}

type Runner interface {
	Run(*CLICommand) error
}

type CLICommand struct {
	Dir          string
	Command      string
	CommandFlags []string
}

type CLI struct {
	ExecCommand string
	Env         []string
	GlobalFlags []string
}

func NewCLI() *CLI {
	return &CLI{
		ExecCommand: ExecCommand,
		Env:         globalEnv,
		GlobalFlags: globalFlags,
	}
}

func (c *CLI) Run(cc *CLICommand) error {
	args := append(append(c.GlobalFlags, cc.Command), cc.CommandFlags...)

	// args are composed from values within this package or directly supplied
	// by the user, who can run arbitrary commands on the system as-is.
	cmd := exec.Command(c.ExecCommand, args...) // #nosec G204

	cmd.Env = c.Env
	cmd.Dir = cc.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c *CLI) Init(dir string) error {
	return c.Run(&CLICommand{
		Dir:          dir,
		Command:      tfInitCommand,
		CommandFlags: tfInitCommandFlags,
	})
}

func (c *CLI) Version() error {
	return c.Run(&CLICommand{
		Command:      tfVersionCommand,
		CommandFlags: tfVersionCommandFlags,
	})
}
