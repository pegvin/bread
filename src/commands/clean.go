package commands

import (
	"os"
	"fmt"
	"bread/src/helpers/utils"
)

type CleanCmd struct {
}

func (cmd *CleanCmd) Run() (err error) {
	// Get the `run-cache` directory path
	appTempDir, err := utils.MakeTempAppDirPath()
	if err != nil {
		return err
	}
	// Remove that directory
	os.RemoveAll(appTempDir)
	fmt.Println("Cache Cleaned")

	registry, err := utils.OpenRegistry()
	if err != nil {
		fmt.Println("Error While Cleaning Registry: " + err.Error())
		fmt.Println("Skipping...")
		return err
	} else {
		registry.Update()
		registry.Close()
		fmt.Println("Registry Cleaned")	
	}

	return nil
}