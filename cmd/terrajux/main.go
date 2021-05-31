package main

import (
	"fmt"
	"os"

	"github.com/rhenning/terrajux"
	"github.com/rhenning/terrajux/internal/cli"
)

// Interfaces:
//   cli.Parser/Runner
//   terraform.Initer/Runner
//   git.Cloner
//   diff.Runner

func main() {
	clii := cli.New(os.Args, terrajux.NewConfig())

	if msg, err := clii.ParseArgs(); err != nil {
		switch err.(type) {

		case *cli.HelpError, *cli.ArgumentError:
			fmt.Println(err.Error())
			fmt.Println(msg)
			os.Exit(1)

		case *cli.VersionError:
			fmt.Print(msg)
			os.Exit(0)
		}
	}
}
