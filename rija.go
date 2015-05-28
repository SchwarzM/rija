package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rija"
	app.Version = Version
	app.Usage = "Gets and sets current issue from JIRA"
	app.Commands = Commands
	app.Run(os.Args)
	app.Action = action
}
