package main

import (
	"github.com/urfave/cli/v2"
)

func command(name string, subcommands ...*cli.Command) (
	c *cli.Command,
) {
	c = &cli.Command{
		Name:        name,
		Subcommands: subcommands,
	}

	return
}

func subcommand(name string, action cli.ActionFunc, flags ...cli.Flag) (
	c *cli.Command,
) {
	c = &cli.Command{
		Name:   name,
		Action: action,
		Flags:  flags,
	}

	return
}

func stringFlag(name string, required bool) (flag *cli.StringFlag) {
	flag = &cli.StringFlag{
		Name:     name,
		Required: required,
	}

	return
}
