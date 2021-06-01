package main

import (
	"fmt"
	"os"

	"github.com/rhenning/terrajux"
	"github.com/rhenning/terrajux/internal/app"
	"github.com/rhenning/terrajux/internal/cli"
)

func main() {
	config := terrajux.NewConfig()
	clii := cli.New(os.Args, config)

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

	appi, err := app.NewDefaultWiring(config)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

	err = appi.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(3)
	}

	os.Exit(0)
}
