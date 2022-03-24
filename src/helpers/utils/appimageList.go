package utils

import (
	"io/ioutil"
	"encoding/json"
	"os/user"
	"path/filepath"
)

// Struct Contains AppImage Item Info
type AppImageFeedItem struct {
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

// Struct Contains Catalog From https://appimage.github.io/feed.jsonhttps://appimage.github.io/feed.json
type AppImageFeed struct {
	Version int
	Home_page_url string
	Feed_url string
	Description string
	Icon string
	Favicon string
	Expired bool
	Items []AppImageFeedItem
}

// Get Full Path to `.AppImageFeed.json` in Applications Dir
func makeAppImageFeedPath() (filePath string, err error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, "Applications", ".AppImageFeed.json"), nil
}

// Read `.AppImageFeed.json` file into a struct
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

// Get the latest information from API
func FetchAppImageListJson() (err error) {
	filePath, err := makeAppImageFeedPath()

	if err != nil {
		return err
	}

	err = DownloadFile("https://appimage.github.io/feed.json", filePath, 0666, "Fetching Json")
	return err
}
