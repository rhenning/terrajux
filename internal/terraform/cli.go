package terraform

import (
	"os"
	"os/exec"
	"strings"
)

const (
	execCommand    = "terraform"
	initCommand    = "init"
	versionCommand = "version"
)

var (
	globalEnv        = []string{"TF_IN_AUTOMATION=1"}
	initCommandFlags = []string{
		"-input=false",
		"-backend=false",
		//"-get-plugins=false",
		//"-upgrade",
	}
)

type cliCommand struct {
	Command      string
	CommandFlags []string
}

// type CLI interface {
// 	Init() (string, string, error)
// 	Run(cliCommand) (string, string, error)
// }

type CLI struct {
	ExecCommand string
	Dir         string
	Env         []string
	GlobalFlags []string
}

func NewCLI(chdir string) *CLI {
	return &CLI{
		ExecCommand: execCommand,
		Dir:         chdir,
		Env:         globalEnv,
		GlobalFlags: []string{},
	}
}

func (c *CLI) Run(cc *cliCommand) error {
	args := append(append(c.GlobalFlags, cc.Command), cc.CommandFlags...)
	cmd := exec.Command(c.ExecCommand, args...)

	for _, v := range c.Env {
		ss := strings.SplitN(v, "=", 2)
		os.Setenv(ss[0], ss[1])
	}

	cmd.Dir = c.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c *CLI) Init() error {
	return c.Run(&cliCommand{
		Command:      initCommand,
		CommandFlags: initCommandFlags,
	})
}

func (c *CLI) Version() error {
	return c.Run(&cliCommand{
		Command:      versionCommand,
		CommandFlags: []string{},
	})
}
