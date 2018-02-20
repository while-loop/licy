package main

import (
	"github.com/urfave/cli"
	"fmt"
)

func init() {
	commands = append(commands, listCommand)
}

var listCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "show all open source licenses",
	Action: func(c *cli.Context) (error) {
		licenses, err := GetLicenses()
		if err != nil {
			fmt.Fprintln(c.App.ErrWriter, err)
			return err
		}

		for _, lic := range licenses {
			fmt.Fprintf(c.App.Writer, "%s (%s)\n", lic.Title, lic.SpdxID)
		}

		return nil
	},
}
