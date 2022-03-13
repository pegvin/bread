package main

import (
	"bread/src/commands"

	"github.com/alecthomas/kong"
)

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Install  commands.InstallCmd  `cmd:"" help:"Install an application."`
	List     commands.ListCmd     `cmd:"" help:"List installed applications."`
	Remove   commands.RemoveCmd   `cmd:"" help:"Remove an application."`
	Update   commands.UpdateCmd   `cmd:"" help:"Update an application."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&commands.Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}