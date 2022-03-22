package repos

import (
	"context"
	"strings"

	"bread/src/helpers/utils"
	"github.com/google/go-github/v31/github"
)

type GitHubRepo struct {
	User    string
	Project string
	Release string
	File    string
}

// Function which parses string to a github repo information, adn returns a object and error (if any)
func NewGitHubRepo(target string) (appInfo Application, err error) {
	appInfo = &GitHubRepo{}

	// Take the `user/repo` and split `user` and `repo`
	targetParts := strings.Split(target, "/")
	ghSource := GitHubRepo{}

	targetPartsLen := len(targetParts)
	if targetPartsLen < 2 { // If input is not in format of `user/repo` assume `user` and `repo` are same
		ghSource = GitHubRepo{
			User:    targetParts[0],
			Project: targetParts[0],
		}
	} else {
		ghSource = GitHubRepo{
			User:    targetParts[0],
			Project: targetParts[1],
		}
	}

	if targetPartsLen > 2 {
		ghSource.Release = targetParts[2]
	}

	if targetPartsLen > 3 {
		ghSource.File = targetParts[3]
	}

	return &ghSource, nil
}

// Function to get the github repo id from the repo information
func (g GitHubRepo) Id() string {
	id := g.User + "/" + g.Project
	return id
}

// Function which gets the latest appimage from github release
func (g GitHubRepo) GetLatestRelease() (*Release, error) {
	var downloadLinks []utils.BinaryUrl // Contains Download Links

	client := github.NewClient(nil) // Client For Interacting with github api
	// Get all the releases from the target
	releases, _, err := client.Repositories.ListReleases(context.Background(), g.User, g.Project, nil)
	if err != nil {
		return nil, err
	}

	// Filter out files which are not AppImage
	for _, release := range releases {
		if *release.Draft {
			continue
		}

		for _, asset := range release.Assets {
			if strings.HasSuffix(*asset.Name, ".AppImage") {
				downloadLinks = append(downloadLinks, utils.BinaryUrl{
					FileName: *asset.Name,
					Url:      *asset.BrowserDownloadURL,
				})
			}
		}

		if len(downloadLinks) > 0 {
			return &Release{
				*release.TagName,
				downloadLinks,
			}, nil
		}
	}

	return nil, NoAppImageBinariesFound
}

// Function which downloads appimage from remote
func (g GitHubRepo) Download(binaryUrl *utils.BinaryUrl, targetPath string) (err error) {
	err = utils.DownloadAppImage(binaryUrl.Url, targetPath)
	return
}

// Function which generates a fallback update information for a appimage
func (g GitHubRepo) FallBackUpdateInfo() string {
	updateInfo := "gh-releases-direct|" + g.User + "|" + g.Project
	if g.Release == "" {
		updateInfo += "|latest"
	} else {
		updateInfo += "|" + g.Release
	}

	if g.File == "" {
		updateInfo += "|*.AppImage"
	} else {
		updateInfo += "|" + g.File
	}

	return updateInfo
}
