package commands

import (
	"os"
	"fmt"
	"bread/src/helpers/utils"
)

type CleanCmd struct {
}

func (cmd *CleanCmd) Run(debug bool) (err error) {
	// Get the `run-cache` directory path
	appTempDir, err := utils.MakeTempAppDirPath()
	if err != nil {
		return err
	}

	// Remove that directory
	os.RemoveAll(appTempDir)
	fmt.Println("Cleaned All The Cache!")
	return nil
}