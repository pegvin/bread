package repos

import (
	"context"
	"bread/src/helpers/utils"
	"strings"

	"github.com/google/go-github/v31/github"
)

type GitHubRepo struct {
	User    string
	Project string
	Release string
	File    string
}

// Function which parses string to a github repo information, adn returns a object and error (if any)
func NewGitHubRepo(target string) (repo Repo, err error) {
	repo = &GitHubRepo{}
	var targetParts []string;

	if strings.HasPrefix(target, "gh:") {
		// Take the `gh:user/repo` Remove "gh:" and split `user` and `repo`
		targetParts = strings.Split(target[3:], "/")
	} else {
		// Take the `gh:user/repo` and split `user` and `repo`
		targetParts = strings.Split(target, "/")
	}

	targetPartsLen := len(targetParts)
	if targetPartsLen < 2 {
		return repo, InvalidTargetFormat
	}

	ghSource := GitHubRepo{
		User:    targetParts[0],
		Project: targetParts[1],
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
	id := "gh:" + g.User + "/" + g.Project

	return id
}

// Function which gets the latest appimage from github release
func (g GitHubRepo) GetLatestRelease() (*Release, error) {
	var downloadLinks []utils.BinaryUrl

	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), g.User, g.Project, nil)
	if err != nil {
		return nil, err
	}

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
