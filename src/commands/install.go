package commands

import (
	"os"
	"fmt"
	"errors"
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"
)

var ApplicationInstalled = errors.New("the application is installed already")

type InstallCmd struct {
	Target string `arg:"" name:"target" help:"Installation target." type:"string"`
	TagName string `arg:"" optional:"" name:"tagname" help:"GitHub Release TagName To Download From." type:"string"`
}

// Function Which Will Be Called When `install` is the Command.
func (cmd *InstallCmd) Run(debug bool) (err error) {
	// Parse The user input
	repo, err := repos.ParseTarget(cmd.Target, cmd.TagName)
	if err != nil {
		return err
	}

	// Get The Latest Release
	release, err := repo.GetLatestRelease()
	if err != nil {
		return err
	}

	if cmd.TagName != "" && release.Tag != cmd.TagName {
		fmt.Println("Tag '" + cmd.TagName + "' not found, using latest available '" + release.Tag + "' instead")
	}

	// Show A Prompt To Select A AppImage File.
	selectedBinary, err := utils.PromptBinarySelection(release.Files)
	if err != nil {
		return err
	}

	// Make A FilePath Out Of The AppImage Name
	targetFilePath, err := utils.MakeTargetFilePath(selectedBinary)
	if err != nil {
		return err
	}

	// Check if the FilePath Exist, Show error
	if _, err = os.Stat(targetFilePath); err == nil {
		return ApplicationInstalled
	}

	// Download The AppImage
	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	// Add The Current Application To The Registry `.registry.json`
	cmd.addToRegistry(targetFilePath, repo, release.Tag, debug)

	// Integrated The AppImage To Desktop
	err = utils.CreateDesktopIntegration(targetFilePath, debug)
	if err != nil {
		fmt.Println("Integration Failed: " + err.Error())
	} else {
		fmt.Println("Integration Complete!")
	}

	// Print Signature Info If Exist.
	utils.ShowSignature(targetFilePath)

	fmt.Println("Installed '" + repo.Id() + "'!")
	return nil
}

// Function To Add Installed Program To Registry (Installed App information is stored in here).
func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo repos.Application, TagName string, debug bool) (error) {
	sha1, _ := utils.GetFileSHA1(targetFilePath) // Get The Sha1 Hash
	updateInfo, _ := utils.ReadUpdateInfo(targetFilePath) // Get The UpdateInfo
	if updateInfo == "" {
		updateInfo = repo.FallBackUpdateInfo()
	}

	appimageInfo, err := utils.GetAppImageInfo(targetFilePath, debug)
	if err != nil {
		return err
	}

	// Make a new entry struct
	entry := utils.RegistryEntry{
		Repo:       repo.Id(),
		TagName:    TagName,
		FileSha1:   sha1,
		// AppName:    "",
		// AppVersion: "",
		FilePath:   targetFilePath,
		UpdateInfo: updateInfo,
		IsTerminalApp: appimageInfo.IsTerminalApp,
		AppImageType: appimageInfo.AppImageType,
	}

	registry, _ := utils.OpenRegistry() // Open The Registry
	if registry != nil {
		_ = registry.Add(entry) // Add the entry to registry `.registry.json`
		_ = registry.Close() // Close the registry
	}
	return nil
}
