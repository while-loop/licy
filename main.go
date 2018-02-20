package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
)

//go:generate go-bindata -o licenses_data.go _licenses

var (
	Name      = "licy"
	Version   = "dev"
	Commit    = "head"
	BuildDate = "now"
	commands  []cli.Command
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = fmt.Sprintf("\nversion:\tv%s\nbuild date:\t%s\ngit hash:\t%s", Version, BuildDate, Commit)
	app.Usage = "open source licensing tool"
	app.Commands = commands
	app.ErrWriter = os.Stderr
	app.Writer = os.Stdout
	app.Run(os.Args)
}
