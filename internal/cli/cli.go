package cli

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/rhenning/terrajux"
)

type ArgumentError struct{}

func (e *ArgumentError) Error() string {
	return "flag: wrong number of arguments"
}

type HelpError struct {
	msg string
}

func (e *HelpError) Error() string {
	return e.msg
}

type VersionError struct{}

func (e *VersionError) Error() string {
	return "flag: version requested"
}

type ArgParser interface {
	ParseArgs() (string, error)
}

type CLI struct {
	Args   []string
	Config *terrajux.Config
}

// New takes a terrajux Config struct and the supplied command-line arguments,
// program-first, just as they are returned from os.Args.
// It returns a CLI.
func New(args []string, config *terrajux.Config) *CLI {
	return &CLI{
		Args:   args,
		Config: config,
	}
}

// ParseArgs processes the supplied command-line arguments and populates the
// appropriate fields in Config from their values.
// It returns an error if the supplied args are invalid.
func (c *CLI) ParseArgs() (message string, err error) {
	var showVersion bool
	var buf bytes.Buffer

	flags := flag.NewFlagSet(c.Args[0], flag.ContinueOnError)

	flags.Usage = func() {
		fmt.Fprintf(
			flags.Output(),
			"\nUsage: %s [options] <giturl> <v1ref> <v2ref> [subpath]\n\nOptions:\n",
			c.Config.Name,
		)

		flags.PrintDefaults()
	}

	flags.SetOutput(&buf)

	flags.BoolVar(&c.Config.CacheClear, "clear", false, "clear cache")
	flags.BoolVar(&showVersion, "version", false, "show version info")
	flags.StringVar(
		&c.Config.DiffTool,
		"difftool",
		"",
		"diff command `template`, e.g. 'opendiff {{.V1}} {{.V2}}'",
	)

	err = flags.Parse(c.Args[1:])

	if err != nil {
		return buf.String(), &HelpError{msg: err.Error()}
	}

	if showVersion {
		flags.Usage = func() {
			fmt.Fprintf(
				flags.Output(),
				"%s %s %s\n",
				c.Config.Name,
				c.Config.Version,
				c.Config.ProjectURL,
			)
		}

		err = &VersionError{}
		flags.Usage()
	} else if flags.NArg() >= 3 && flags.NArg() <= 4 {
		c.Config.GitURL = flags.Arg(0)
		c.Config.GitRefV1 = flags.Arg(1)
		c.Config.GitRefV2 = flags.Arg(2)
		c.Config.GitSubpath = flags.Arg(3)
	} else {
		err = &ArgumentError{}
		flags.Usage()
	}

	message = buf.String()

	return message, err
}
