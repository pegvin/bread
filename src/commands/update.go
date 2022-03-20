package commands

import (
	"errors"
	"fmt"
	"strings"

	"bread/src/helpers/utils"
	update "github.com/DEVLOPRR/appimage-update"
	"github.com/DEVLOPRR/appimage-update/updaters"
)

type UpdateCmd struct {
	Targets []string `arg:"" optional:"" name:"targets" help:"Updates the target applications." type:"string"`

	Check bool `help:"Only check for updates."`
	All   bool `help:"Update all applications."`
}

var NoUpdateInfo = errors.New("there is no update information")

// Function Which Will Be Executed When `update` is called.
func (cmd *UpdateCmd) Run() (err error) {
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

		updateMethod, err := NewUpdater(entry.UpdateInfo, entry.FilePath)
		if err != nil {
			println(err.Error())
			continue
		}

		fmt.Println("Looking for updates of: ", entry.FilePath)
		updateAvailable, err := updateMethod.Lookup()
		if err != nil {
			println(err.Error())
			continue
		}

		if !updateAvailable {
			fmt.Println("No updates were found for: ", entry.FilePath)
			continue
		}

		if cmd.Check {
			fmt.Println("Update available for: ", entry.FilePath)
			continue
		}

		result, err := updateMethod.Download()
		if err != nil {
			println(err.Error())
			continue
		}

		utils.ShowSignature(result)
		fmt.Println("Update downloaded to: " + result)
	}

	return nil
}

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

func NewUpdater(updateInfoString string, appImagePath string) (update.Updater, error) {
	if strings.HasPrefix(updateInfoString, "zsync") {
		return updaters.NewZSyncUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "gh-releases-zsync") {
		return updaters.NewGitHubZsyncUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "gh-releases-direct") {
		return updaters.NewGitHubDirectUpdater(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "ocs-v1-appimagehub-direct") {
		return updaters.NewOCSAppImageHubDirect(&updateInfoString, appImagePath)
	}

	if strings.HasPrefix(updateInfoString, "ocs-v1-appimagehub-zsync") {
		return updaters.NewOCSAppImageHubZSync(&updateInfoString, appImagePath)
	}

	return nil, fmt.Errorf("invalid updated information: %q", updateInfoString)
}
