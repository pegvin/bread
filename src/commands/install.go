package commands

import (
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"

	"errors"
	"fmt"
	"os"
)

var ApplicationInstalled = errors.New("the application is installed already")

type InstallCmd struct {
	Target string `arg:"" name:"target" help:"Installation target." type:"string"`
}

// Function Which Will Be Called When `install` is the Command.
func (cmd *InstallCmd) Run() (err error) {
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
	cmd.addToRegistry(targetFilePath, repo)

	// Integrated The AppImage To Desktop
	cmd.createDesktopIntegration(targetFilePath)

	// Print Signature Info If Exist.
	signingEntity, _ := utils.VerifySignature(targetFilePath)
	if signingEntity != nil {
		fmt.Println("AppImage signed by:")
		for _, v := range signingEntity.Identities {
			fmt.Println("\t", v.Name)
		}
	}

	return
}

// Function To Add Installed Program To Registry (Installed App information is stored in here).
func (cmd *InstallCmd) addToRegistry(targetFilePath string, repo repos.Application) {
	sha1, _ := utils.GetFileSHA1(targetFilePath) // Get The Sha1 Hash
	updateInfo, _ := utils.ReadUpdateInfo(targetFilePath) // Get The UpdateInfo
	if updateInfo == "" {
		updateInfo = repo.FallBackUpdateInfo()
	}

	// Make a new entry struct
	entry := utils.RegistryEntry{
		FilePath:   targetFilePath,
		Repo:       repo.Id(),
		FileSha1:   sha1,
		UpdateInfo: updateInfo,
	}

	registry, _ := utils.OpenRegistry() // Open The Registry
	if registry != nil {
		_ = registry.Add(entry) // Add the entry to registry `.registry.json`
		_ = registry.Close() // Close the registry
	}
}

// Function To Integrate The AppImage To Desktop. (Can Only Be Called From InstallCmd Struct)
func (cmd *InstallCmd) createDesktopIntegration(targetFilePath string) {
	libAppImage, err := utils.NewLibAppImageBindings() // Load the `libappimage` Library For Integration
	if err != nil {
		fmt.Println("Integration failed:", err.Error())
		return
	}

	err = libAppImage.Register(targetFilePath) // Register The File
	if err != nil {
		fmt.Println("Integration failed: " + err.Error())
	} else {
		fmt.Println("Integration completed")
	}
}
