package main

import (
	"fmt"
	"bread/src/commands"
	"github.com/alecthomas/kong"
)

// Variable Which will be set on the compile time using ldflags
var VERSION string

type VersionFlag bool

var cli struct {
	Install    commands.InstallCmd    `cmd:"" help:"Install an application."`
	Run        commands.RunCmd        `cmd:"" help:"Run an application from Remote."`
	List       commands.ListCmd       `cmd:"" help:"List installed applications."`
	Remove     commands.RemoveCmd     `cmd:"" help:"Remove an application."`
	Update     commands.UpdateCmd     `cmd:"" help:"Update an application."`
	Search     commands.SearchCmd     `cmd:"" help:"Search for appliation from appimage list."`
	Clean      commands.CleanCmd      `cmd:"" help:"Clean all the cache & unused registry entries."`
	Version    VersionFlag            `name:"version" short:"v" help:"Print version information and quit"`
	Debug      bool                   `short:"d" help:"Show extra information for debugging." default:"false"`
}

func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	if VERSION == "" {
		fmt.Println("Unknown Custom Build")
	} else {
		fmt.Println("Bread v" + VERSION)
	}
	app.Exit(0)
	return nil
}

func main() {
	ctx := kong.Parse(
		&cli,
		kong.Name("bread"),
		kong.Description("Install, update, remove & run AppImage from GitHub using your CLI."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}))
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(cli.Debug)
	ctx.FatalIfErrorf(err)
}
