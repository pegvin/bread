package commands

import (
	"os"
	"fmt"
	"bread/src/helpers/utils"
)

type CleanCmd struct {
}

func (cmd *CleanCmd) Run(debug bool) (err error) {
	appTempDir, err := utils.MakeTempAppDirPath()
	if err != nil {
		return err
	}

	os.RemoveAll(appTempDir)
	fmt.Println("Cleaned All The Cache!")
	return nil
}