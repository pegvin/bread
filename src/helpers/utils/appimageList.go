package utils

import (
	"io/ioutil"
	"encoding/json"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
)

type AppImageFeedItems struct {
	Name string
	Description string
	Categories []string
	Authors []struct {
		Name string
		Url string
	}
	License string
	Links []struct {
		Type string
		Url string
	}
	Icons []string
	Screenshots []string
}

type AppImageFeed struct {
	Version int
	Home_page_url string
	Feed_url string
	Description string
	Icon string
	Favicon string
	Expired bool
	Items []AppImageFeedItems
}

func makeAppImageFeedPath() (filePath string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, "Applications", ".AppImageFeed.json"), nil
}

func ReadAppImageListJson() (aifeedJson *AppImageFeed, err error) {
	filePath, err := makeAppImageFeedPath()

	if err != nil {
		return nil, err
	}

	myAppImageJson := &AppImageFeed{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &myAppImageJson)
	if err != nil {
		return nil, err
	}
	return myAppImageJson, nil
}

func FetchAppImageListJson() (err error) {
	filePath, err := makeAppImageFeedPath()

	if err != nil {
		return err
	}

	output, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer output.Close()

	resp, err := http.Get("https://appimage.github.io/feed.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Fetching Json",
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
