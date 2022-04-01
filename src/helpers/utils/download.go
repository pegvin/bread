package utils

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"github.com/schollz/progressbar/v3"
)

// Download a file from remote
func DownloadFile(url string, filePath string, permission fs.FileMode, barText string) (err error) {
	appDir, err := MakeApplicationsDirPath()
	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir(appDir, "temp")
	if err != nil {
		return err
	}

	tempFilePath := tempDir + "/" + filepath.Base(filePath)
	output, err := os.OpenFile(tempFilePath, os.O_RDWR|os.O_CREATE, permission)
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

	// Handles Ctrl + C Detection
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		_ = resp.Body.Close()
		_ = output.Close()
		_ = os.Remove(filePath)
		_ = os.RemoveAll(tempDir)

		fmt.Println("Ctrl + C, Removing Downloaded File & Exiting.")
		os.Exit(0)
	}()

	_, err = io.Copy(io.MultiWriter(output, bar), resp.Body)
	err = os.Rename(tempFilePath, filePath)
	if err != nil {
		return err
	}

	return os.RemoveAll(tempDir)
}
