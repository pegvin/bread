package commands

import (
	"os"
	"fmt"
	"sort"
	"bytes"
	"path/filepath"
	"bread/src/helpers/utils"
	"github.com/juju/ansiterm"
)

type ListCmd struct {
	ShowSha1 bool `short:"s" name:"show-sha1" help:"Show SHA1 Hashes too." default:"false"`
}

// Function which will be executed when `list` is called.
func (r *ListCmd) Run() error {
	registry, err := utils.OpenRegistry() // Open The Registry
	if err != nil {
		return err
	}
	defer registry.Close() // Close the registry before function return

	registry.Update() // Update the registry with latest application info
	if len(registry.Entries) == 0 {
		fmt.Println("No installed Applications Found!")
		return nil
	}

	var buf bytes.Buffer // Buffer which will hold the table
	tabWriter := ansiterm.NewTabWriter(&buf, 20, 4, 0, ' ', 0)
	tabWriter.SetColorCapable(true)

	tabWriter.SetForeground(ansiterm.Green)
	if r.ShowSha1 {
		_, _ = tabWriter.Write([]byte("User/Repo\t File Name\t SHA1 HASH\n"))
		_, _ = tabWriter.Write([]byte("---------\t ---------\t ---------\n"))
	} else {
		_, _ = tabWriter.Write([]byte("User/Repo\t File Name\n"))
		_, _ = tabWriter.Write([]byte("---------\t ---------\n"))
	}

	tabWriter.SetForeground(ansiterm.DarkGray)

	var lines [][]string
	for fileName, v := range registry.Entries {
		var line []string
		if r.ShowSha1 {
			line = []string{v.Repo, filepath.Base(fileName), v.FileSha1}
		} else {
			line = []string{v.Repo, filepath.Base(fileName)}
		}
		lines = append(lines, line)
	}
	sort.Slice(lines, func(i int, j int) bool {
		return lines[i][1] < lines[j][1]
	})

	for _, line := range lines {
		if r.ShowSha1 {
			_, _ = tabWriter.Write([]byte(line[0] + "\t " + line[1] + "\t " + line[2] + "\n"))
		} else {
			_, _ = tabWriter.Write([]byte(line[0] + "\t " + line[1] + "\t " + "\n"))
		}
	}
	_ = tabWriter.Flush()
	_, _ = os.Stdout.Write(buf.Bytes())
	return nil
}
