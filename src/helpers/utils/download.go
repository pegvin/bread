package utils

import (
	"io"
	"os"
	"io/fs"
	"net/http"
	"os/signal"
	"github.com/schollz/progressbar/v3"
)

// Download a file from remote
func DownloadFile(url string, filePath string, permission fs.FileMode, barText string) (err error) {
	output, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, permission)
	if err != nil {
		return err
	}
	defer output.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		barText,
	)

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		_ = resp.Body.Close()
		_ = output.Close()
		_ = os.Remove(filePath)

		os.Exit(0)
	}()

	_, err = io.Copy(io.MultiWriter(output, bar), resp.Body)
	return err
}
