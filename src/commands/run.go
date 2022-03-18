package commands

import (
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"

	"os"
	"os/exec"
)

type RunCmd struct {
	Target string `arg:"" name:"target" help:"Target To Run" type:"string"`
}

func (cmd *RunCmd) Run() (err error) {
	// Parse The user input
	repo, err := repos.ParseTarget(cmd.Target)
	if err != nil {
		return err
	}

	// Get The Latest Release
	release, err := repo.GetLatestRelease()
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

	// Check if the FilePath Exist, Show error
	if _, err = os.Stat(targetFilePath); err == nil {
		command := exec.Command(targetFilePath)
		command.Run()
		return nil
	}

	// Download The AppImage
	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	// Print Signature Info If Exist.
	utils.ShowSignature(targetFilePath)

	command := exec.Command(targetFilePath)
	return command.Run()
}