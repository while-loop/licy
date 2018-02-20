package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"time"
)

func init() {
	commands = append(commands, genCommand)
}

var genCommand = cli.Command{
	Name:    "gen",
	Aliases: []string{"get"},
	Usage:   "generate and populate an open source license",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "p",
			Usage: "project name",
		}, cli.StringFlag{
			Name:  "o",
			Usage: "output filename",
			Value: "LICENSE",
		}, cli.IntFlag{
			Name:  "y",
			Usage: "year of license init",
			Value: time.Now().Year(),
		},
	},
	Action: func(c *cli.Context) error {
		license := c.Args().Get(0)
		lic, err := GetLicense(license)
		if err == ErrNotFound {
			suggestMispell(c.App.ErrWriter, license)
			return err
		} else if err != nil {
			fmt.Fprintln(c.App.ErrWriter, err)
			return err
		}

		name := ""
		if c.NArg() >= 2 {
			name = c.Args().Get(1)
		}
		lic.FillBody(name, c.Int("y"))
		if c.String("p") != "" {
			lic.FillProject(c.String("p"))
		}

		f, err := os.Create(c.String("o"))
		if err != nil {
			fmt.Fprintf(c.App.ErrWriter, "Unable to create output license file %s: %v\n", c.String("o"), err)
			return err
		}
		defer f.Close()

		if _, err = lic.Save(f); err != nil {
			fmt.Fprintf(c.App.ErrWriter, "Unable write to license file: %v\n", err)
			return err
		}

		return nil
	},
}
