package commands

import (
	"os"
	"errors"
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"
	"github.com/mgord9518/aisap"
)

type RunCmd struct {
	Target         string `arg:"" name:"target" help:"Target To Run" type:"string"`
	Level          uint8   `arg:"" help:"Set Permission Level" type:"int" default:"0"`
	Arguments    []string `arg:"" passthrough:"" optional:"" name:"arguments" help:"Argument to pass to the program" type:"string"`
	NoPreRelease   bool `short:"n" help:"Disable pre-releases." default:"false"`
}

func runAppImage(filePath string, permissionLevel uint8, arguments []string) (error) {
	appImage, err := aisap.NewAppImage(filePath)
	if err != nil {
		return nil
	}
	err = appImage.Perms.SetLevel(int(permissionLevel))
	if err != nil {
		return err
	}
	return appImage.Run(arguments)
}

func (cmd *RunCmd) Run(debug bool) (err error) {
	if cmd.Level > 4 {
		return errors.New("permission level can only be 0, 1, 2 or 3")
	}
	// Parse The user input
	repo, err := repos.ParseTarget(cmd.Target, "")
	if err != nil {
		return err
	}

	// Get The Latest Release
	release, err := repo.GetLatestRelease(cmd.NoPreRelease)
	if err != nil {
		return err
	}

	// Show A Prompt To Select A AppImage File.
	selectedBinary, err := utils.PromptBinarySelection(release.Files)
	if err != nil {
		return err
	}

	// Make A FilePath Out Of The AppImage Name
	targetFilePath, err := utils.MakeTempFilePath(selectedBinary)
	if err != nil {
		return err
	}

	// Check if the FilePath Exist (cached file), Show error
	if _, err = os.Stat(targetFilePath); err == nil {
		return runAppImage(targetFilePath, cmd.Level, cmd.Arguments)
	}

	// Download The AppImage
	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	// Print Signature Info If Exist.
	utils.ShowSignature(targetFilePath)

	return runAppImage(targetFilePath, cmd.Level, cmd.Arguments)
}