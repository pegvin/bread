package commands

import (
	"fmt"
	"bread/src/helpers/utils"
	"github.com/DEVLOPRR/libappimage-go"
	"os"
	"strings"
)

type RemoveCmd struct {
	Target   string `arg:"" name:"target" help:"target to remove" type:"string"`
	KeepFile bool   `help:"Remove only the application desktop entry."`
}

// Function which will be executed when `remove` is called.
func (cmd *RemoveCmd) Run(debug bool) (err error) {
	cmd.Target = strings.ToLower(cmd.Target)
	registry, err := utils.OpenRegistry() // Open The Registry
	if err != nil {
		return err
	}
	defer registry.Close() // Close the registry before function end

	registry.Update() // Update the registry with latest installed appimage info

	// If the user provided string is short like `libresprite` convert it to `libresprite/libresprite`
	if len(strings.Split(cmd.Target, "/")) < 2 {
		cmd.Target = cmd.Target + "/" + cmd.Target;
	}

	entry, ok := registry.Lookup(cmd.Target) // Find the application in the registry
	if !ok {
		return fmt.Errorf("application not found \"" + cmd.Target + "\"")
	}

	err = removeDesktopIntegration(entry.FilePath, debug) // Remove the application desktop integration
	if err != nil {
		fmt.Println("Desktop integration removal failed: " + err.Error())
	}

	if cmd.KeepFile {
		return nil
	}

	err = os.Remove(entry.FilePath)
	if err != nil {
		return fmt.Errorf("Unable to remove AppImage file: " + err.Error())
	}
	fmt.Println("Application removed: " + entry.FilePath)
	registry.Remove(entry.FilePath)
	return err
}

// Remove the application desktop integration
func removeDesktopIntegration(filePath string, debug bool) error {
	libAppImage, err := libappimagego.NewLibAppImageBindings()
	if err != nil {
		return err
	}

	if libAppImage.ShallNotBeIntegrated(filePath) {
		return nil
	}

	return libAppImage.Unregister(filePath, debug)
}
