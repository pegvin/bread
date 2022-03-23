package main

import (
	"fmt"
	"bread/src/commands"
	"github.com/alecthomas/kong"
)

type VersionFlag bool

var cli struct {
	Install    commands.InstallCmd    `cmd:"" help:"Install an application."`
	Run        commands.RunCmd        `cmd:"" help:"Run an application from Remote."`
	List       commands.ListCmd       `cmd:"" help:"List installed applications."`
	Remove     commands.RemoveCmd     `cmd:"" help:"Remove an application."`
	Update     commands.UpdateCmd     `cmd:"" help:"Update an application."`
	Search     commands.SearchCmd     `cmd:"" help:"Search for appliation from appimage list."`
	Clean      commands.CleanCmd      `cmd:"" help:"Clean all the cache."`
	Version    VersionFlag            `name:"version" help:"Print version information and quit"`
	Debug      bool                   `help:"Show extra information for debugging."`
}

func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Bread v" + vars["VERSION"])
	app.Exit(0)
	return nil
}

func main() {
	ctx := kong.Parse(
		&cli,
		kong.Name("bread"),
		kong.Description("Install, update and remove AppImage from GitHub using your CLI."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"VERSION": "0.4.4",
		})
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(cli.Debug)
	ctx.FatalIfErrorf(err)
}
