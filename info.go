package main

import (
	"github.com/urfave/cli"
	"fmt"
	"text/tabwriter"
)

func init() {
	commands = append(commands, infoCommand)
}

var infoCommand = cli.Command{
	Name:    "info",
	Aliases: []string{"show", "more"},
	Usage:   "show information about a given license",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "v",
			Usage: "verbose",
		},
	},
	Action: func(c *cli.Context) error {
		given := c.Args().First()
		lic, err := GetLicense(given)
		if err == ErrNotFound {
			suggestMispell(c.App.ErrWriter, given)
			return err
		} else if err != nil {
			fmt.Fprintln(c.App.ErrWriter, err)
			return err
		}

		fmt.Fprintf(c.App.Writer, "%s (%s)\n", lic.Title, lic.SpdxID)
		fmt.Fprintln(c.App.Writer, lic.Source)
		fmt.Fprintln(c.App.Writer, lic.Description)
		fmt.Fprintln(c.App.Writer)
		w := tabwriter.NewWriter(c.App.Writer, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "permissions\tconditions\tlimitations")
		fmt.Fprintln(w, "-----------\t----------\t-----------")
		for i := 0; i < len(lic.Permissions) || i < len(lic.Conditions) || i < len(lic.Limitations); i++ {
			if i < len(lic.Permissions) {
				fmt.Fprint(w, lic.Permissions[i])
			}
			fmt.Fprint(w, "\t")
			if i < len(lic.Conditions) {
				fmt.Fprint(w, lic.Conditions[i])
			}
			fmt.Fprint(w, "\t")
			if i < len(lic.Limitations) {
				fmt.Fprint(w, lic.Limitations[i])
			}
			fmt.Fprintln(w, "\t")

		}
		w.Flush()
		if c.Bool("v") {
			fmt.Fprintln(c.App.Writer, lic.Body)
		}

		return nil
	},
}
