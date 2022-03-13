package commands

import (
	"bytes"
	"bread/src/helpers/utils"
	"os"
	"path/filepath"
	"sort"

	"github.com/juju/ansiterm"
)

type ListCmd struct {
}

// Function which will be executed when `list` is called.
func (r *ListCmd) Run(*Context) error {
	registry, err := utils.OpenRegistry()
	if err != nil {
		return err
	}
	defer registry.Close()

	registry.Update()
	var buf bytes.Buffer
	tabWriter := ansiterm.NewTabWriter(&buf, 20, 4, 0, ' ', 0)
	tabWriter.SetColorCapable(true)

	tabWriter.SetForeground(ansiterm.Green)
	_, _ = tabWriter.Write([]byte("Host\t File Name\t SHA1\n"))
	_, _ = tabWriter.Write([]byte("----\t ---------\t ----\n"))

	tabWriter.SetForeground(ansiterm.DarkGray)

	var lines [][]string
	for fileName, v := range registry.Entries {
		line := []string{v.Repo, filepath.Base(fileName), v.FileSha1}
		lines = append(lines, line)
	}
	sort.Slice(lines, func(i int, j int) bool {
		return lines[i][1] < lines[j][1]
	})

	for _, line := range lines {
		_, _ = tabWriter.Write([]byte(line[0] + "\t " + line[1] + "\t " + line[2] + "\n"))
	}
	_ = tabWriter.Flush()
	_, _ = os.Stdout.Write(buf.Bytes())
	return nil
}
