package commands

import (
	"bread/src/helpers/repos"
	"bread/src/helpers/utils"

	"fmt"
	"os"
	goCmd "github.com/go-cmd/cmd"
)

type RunCmd struct {
	Target string `arg:"" name:"target" help:"Target To Run" type:"string"`
	Arguments []string `arg:"" passthrough:"" optional:"" name:"arguments" help:"Argument to pass to the program" type:"string"`
}

func executeCmd(target string, arguments []string) {
	options := goCmd.Options{
		Buffered: false,
		Streaming: true,
	}
	runCmd := goCmd.NewCmdOptions(options, target, arguments...)

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for runCmd.Stdout != nil || runCmd.Stderr != nil {
			select {
			case line, open := <-runCmd.Stdout:
				if !open {
					runCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-runCmd.Stderr:
				if !open {
					runCmd.Stderr = nil
					continue
				}
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-runCmd.Start()

	// Wait for goroutine to print everything
	<-doneChan
}

func (cmd *RunCmd) Run() (err error) {
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
	targetFilePath, err := utils.MakeTempFilePath(selectedBinary)
	if err != nil {
		return err
	}

	// Check if the FilePath Exist, Show error
	if _, err = os.Stat(targetFilePath); err == nil {
		executeCmd(targetFilePath, cmd.Arguments)
		return nil
	}

	// Download The AppImage
	err = repo.Download(selectedBinary, targetFilePath)
	if err != nil {
		return err
	}

	// Print Signature Info If Exist.
	utils.ShowSignature(targetFilePath)

	executeCmd(targetFilePath, cmd.Arguments)
	return nil
}