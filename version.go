package main

import (
	"github.com/urfave/cli"
)

func init() {
	commands = append(commands, versionCommand)
}

var versionCommand = cli.Command{
	Name:    "version",
	Usage:   "licy",
	Action: func(c *cli.Context) (error) {
		cli.ShowVersion(c)
		return nil
	},
}
