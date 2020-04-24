package main

import (
	"log"
	"os"

	"github.com/quarantine_cli/cmd"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	info(app)
	cmd.New(app)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func info(app *cli.App) {
	app.Name = "The quarantine CLI"
	app.Usage = "Come here to get all your quarantine info, or at least what I have added..."
	app.Author = "knee-knee"
}
