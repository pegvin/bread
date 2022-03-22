package commands

import (
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"
	"github.com/DEVLOPRR/libappimage-go"

	"errors"
	"fmt"
	"os"
)

var ApplicationInstalled = errors.New("the application is installed already")

type InstallCmd struct {
	Target string `arg:"" name:"target" help:"Installation target." type:"string"`
}

// Function Which Will Be Called When `install` is the Command.
func (cmd *InstallCmd) Run(debug bool) (err error) {
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
	cmd.addToRegistry(targetFilePath, repo, debug)

	// Integrated The AppImage To Desktop
	cmd.createDesktopIntegration(targetFilePath, debug)

	// Print Signature Info If Exist.
	utils.ShowSignature(targetFilePath)

	return nil
}

// Function To Add Installed Program To Registry (Installed App information is stored in here).
func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo repos.Application, debug bool) (error) {
	sha1, _ := utils.GetFileSHA1(targetFilePath) // Get The Sha1 Hash
	updateInfo, _ := utils.ReadUpdateInfo(targetFilePath) // Get The UpdateInfo
	if updateInfo == "" {
		updateInfo = repo.FallBackUpdateInfo()
	}

	appimageInfo, err := getAppImageInfo(targetFilePath, debug)
	if err != nil {
		return err
	}

	// Make a new entry struct
	entry := utils.RegistryEntry{
		Repo:       repo.Id(),
		FileSha1:   sha1,
		AppName:    "",
		AppVersion: "",
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

func getAppImageInfo(targetFilePath string, debug bool) (*utils.AppImageInfo, error) {
	libAppImage, err := libappimagego.NewLibAppImageBindings() // Load the `libappimage` Library For Integration
	if err != nil {
		return nil, err
	}

	return &utils.AppImageInfo{
		IsTerminalApp: libAppImage.IsTerminalApp(targetFilePath),
		AppImageType: libAppImage.GetType(targetFilePath, debug),
	}, nil
}

// Function To Integrate The AppImage To Desktop. (Can Only Be Called From InstallCmd Struct)
func (cmd *InstallCmd) createDesktopIntegration(targetFilePath string, debug bool) {
	libAppImage, err := libappimagego.NewLibAppImageBindings() // Load the `libappimage` Library For Integration
	if err != nil {
		fmt.Println("Integration failed:", err.Error())
		return
	}

	err = libAppImage.Register(targetFilePath, debug) // Register The File
	if err != nil {
		fmt.Println("Integration failed: " + err.Error())
	} else {
		fmt.Println("Integration completed")
	}
}
