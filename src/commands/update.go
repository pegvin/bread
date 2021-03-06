package commands

import (
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UpdateCmd struct {
	Targets []string `arg:"" optional:"" name:"targets" help:"Update the target/all applications." type:"string"`

	Check          bool `short:"c" help:"Only check for updates."`
	All            bool `short:"a" help:"Update/check all applications."`
	NoPreRelease   bool `short:"n" help:"Disable pre-releases." default:"false"`
}

// Function Which Will Be Executed When `update` is called.
func (cmd *UpdateCmd) Run() (err error) {
	// Variable which will hold if any app was updated.
	var howManyUpdates int = 0

	if cmd.All { // if `update all`
		cmd.Targets, err = getAllTargets() // Load all the application info into targets
		if err != nil {
			return err
		}
	}

	if len(cmd.Targets) == 0 {
		if cmd.All {
			fmt.Println("No Application Installed")
		} else {
			fmt.Println("No Application Specified To Update")
		}
		return nil
	}

	fmt.Println("Checking For Updates")
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
	
		release, err := repo.GetLatestRelease(cmd.NoPreRelease)
		if err != nil {
			fmt.Println(target, "\U00002192", err)
			continue
			// return err
		}

		if release.Tag == entry.TagName {
			continue
		}

		if cmd.Check {
			fmt.Println(target, "\U00002192", release.Tag)
			howManyUpdates++
			continue
		}

		fmt.Println("Updating: " + target + "#" + entry.TagName + " \U00002192 " + target + "#" + release.Tag)

		var selectedBinary *utils.BinaryUrl;
		for fileIndex := range release.Files {
			if filepath.Base(entry.FilePath) == release.Files[fileIndex].FileName {
				selectedBinary = &release.Files[fileIndex]
				break
			}
		}

		if selectedBinary == nil {
			// Show A Prompt To Select A AppImage File.
			selectedBinary, err = utils.PromptBinarySelection(release.Files)
			if err != nil {
				return err
			}
		}

		// Make A FilePath Out Of The AppImage Name
		targetFilePath, err := utils.MakeTargetFilePath(selectedBinary)
		if err != nil {
			return err
		}

		// Download The AppImage
		err = repo.Download(selectedBinary, targetFilePath)
		if err != nil {
			return err
		}

		// Integrated The AppImage To Desktop
		err = utils.CreateDesktopIntegration(targetFilePath)
		if err != nil {
			os.Remove(targetFilePath)
			return err
		}

		registry, err := utils.OpenRegistry()
		if err != nil {
			return err
		}

		// De-Integrate old app from desktop
		err = utils.RemoveDesktopIntegration(entry.FilePath)
		if err != nil {
			os.Remove(targetFilePath) // If error, remove the newly downloaded appimage.
			return err
		} else {
			err = os.Remove(entry.FilePath) // Remove the old appimage
			if err != nil {
				fmt.Println("Cannot Remove The Old AppImage.\n", err.Error())
			}
		}

		registry.Remove(entry.FilePath) // Remove Old File From Registry

		sha1hash, _ := utils.GetFileSHA1(targetFilePath)
		appImageInfo, _ := utils.GetAppImageInfo(targetFilePath)
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

		err = registry.Close()
		if err != nil {
			return err
		}

		// Print Signature Info If Exist.
		utils.ShowSignature(targetFilePath)

		// Remove the old file
		os.Remove(entry.FilePath)

		// utils.ShowSignature(result)
		fmt.Println("Updated: " + target)
		howManyUpdates++
	}

	if cmd.Check {
		if howManyUpdates == 0 {
			fmt.Println("No Updates Found!")
		} else {
			fmt.Println("Update Available For", howManyUpdates, "Application(s)")
		}
	} else {
		if howManyUpdates == 0 {
			fmt.Println("No Updates Found!")
		} else {
			fmt.Println("Updated", howManyUpdates, "Application(s)")
		}
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

	return entry, nil
}

// Get all the applications from the registry
func getAllTargets() ([]string, error) {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return nil, err
	}
	registry.Update()

	var repos []string
	for k := range registry.Entries {
		entry, _ := registry.Lookup(k)
		repos = append(repos, entry.Repo)
	}

	return repos, nil
}
