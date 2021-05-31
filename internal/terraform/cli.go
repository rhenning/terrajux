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
	Init() error
}

type Runner interface {
	Run(*CLICommand) error
}

type CLICommand struct {
	Command      string
	CommandFlags []string
}

type CLI struct {
	ExecCommand string
	Dir         string
	Env         []string
	GlobalFlags []string
}

func NewCLI(chdir string) *CLI {
	return &CLI{
		ExecCommand: ExecCommand,
		Dir:         chdir,
		Env:         globalEnv,
		GlobalFlags: globalFlags,
	}
}

func (c *CLI) Run(cc *CLICommand) error {
	args := append(append(c.GlobalFlags, cc.Command), cc.CommandFlags...)
	cmd := exec.Command(c.ExecCommand, args...)

	cmd.Env = append(os.Environ(), c.Env...)
	cmd.Dir = c.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c *CLI) Init() error {
	return c.Run(&CLICommand{
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
