package commands

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"bread/src/helpers/utils"
	"bread/src/helpers/repos"
)

type UpdateCmd struct {
	Targets []string `arg:"" optional:"" name:"targets" help:"Update the target/all applications." type:"string"`

	Check bool `help:"Only check for updates."`
	All   bool `help:"Update all applications."`
}

var NoUpdateInfo = errors.New("there is no update information")

// Function Which Will Be Executed When `update` is called.
func (cmd *UpdateCmd) Run(debug bool) (err error) {
	fmt.Println("Checking For Updates")

	if cmd.All { // if `update all`
		cmd.Targets, err = getAllTargets() // Load all the application info into targets
		if err != nil {
			return err
		}
	}

	for _, target := range cmd.Targets {
		if len(strings.Split(target, "/")) < 2 {
			target = strings.ToLower(target + "/" + target)
		} else if len(strings.Split(target, "/")) == 2 {
			target = strings.ToLower(target)
		}
		entry, err := cmd.getRegistryEntry(target)
		if err != nil {
			continue
		}

		repo, err := repos.ParseTarget(target, "")

		if err != nil {
			return err
		}
	
		release, err := repo.GetLatestRelease()
		if err != nil {
			return err
		}
	
		if release.Tag != entry.TagName {
			if cmd.Check {
				fmt.Println("Update Available: " + target + "#" + release.Tag)
				return nil
			}

			fmt.Println("Updating: " + target + "#" + entry.TagName + " \U00002192 " + target + "#" + release.Tag)

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

			registry, err := utils.OpenRegistry()
			registry.Remove(entry.FilePath) // Remove old file from registry

			// // Integrated The AppImage To Desktop
			err = utils.CreateDesktopIntegration(targetFilePath, debug)
			if err != nil {
				os.Remove(targetFilePath)
				return err
			}

			sha1hash, _ := utils.GetFileSHA1(targetFilePath)
			appImageInfo, _ := utils.GetAppImageInfo(targetFilePath, debug)
			err = registry.Add(utils.RegistryEntry{
				Repo: target,
				FilePath: targetFilePath,
				FileSha1: sha1hash,
				TagName: release.Tag,
				IsTerminalApp: appImageInfo.IsTerminalApp,
				AppImageType: appImageInfo.AppImageType,
			})

			if err != nil {
				return err
			}

			// De-Integrate old app from desktop
			err = utils.RemoveDesktopIntegration(entry.FilePath, debug)
			if err != nil {
				os.Remove(targetFilePath)
				return err
			}

			registry.Remove(entry.FilePath)
			err = registry.Close()
			if err != nil {
				return err
			}

			// Print Signature Info If Exist.
			utils.ShowSignature(targetFilePath)

			// Remove the old file
			os.Remove(entry.FilePath)
		}

		// utils.ShowSignature(result)
		fmt.Println("Updated: " + target)
	}

	return nil
}

// Get a application from registry
func (cmd *UpdateCmd) getRegistryEntry(target string) (utils.RegistryEntry, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return utils.RegistryEntry{}, err
	}
	defer registry.Close()

	entry, _ := registry.Lookup(target)

	if entry.UpdateInfo == "" {
		entry.UpdateInfo, _ = utils.ReadUpdateInfo(target)
		entry.FilePath = target
	}

	if entry.UpdateInfo == "" {
		return entry, NoUpdateInfo
	} else {
		return entry, nil
	}
}

// Get all the applications from the registry
func getAllTargets() ([]string, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return nil, err
	}
	registry.Update()

	var paths []string
	for k := range registry.Entries {
		paths = append(paths, k)
	}

	return paths, nil
}
