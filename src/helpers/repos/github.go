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
	TagName string
}

// Function which parses string to a github repo information, adn returns a object and error (if any)
func NewGitHubRepo(target string, tagName string) (appInfo Application, err error) {
	appInfo = &GitHubRepo{}
	ghSource := GitHubRepo{}

	userRepo, err := utils.GetUserRepoFromUrl(target)
	if err == nil {
		userRepoSplitted := strings.Split(userRepo, "/")
		ghSource = GitHubRepo{
			User:    userRepoSplitted[0],
			Project: userRepoSplitted[1],
			TagName: tagName,
		}
		return &ghSource, nil
	} else {
		// Take the `user/repo` and split `user` and `repo`
		targetParts := strings.Split(target, "/")

		// If input is not in format of `user/repo` assume `user` and `repo` are same
		if len(targetParts) < 2 {
			ghSource = GitHubRepo{
				User:    targetParts[0],
				Project: targetParts[0],
				TagName: tagName,
			}
		} else {
			ghSource = GitHubRepo{
				User:    targetParts[0],
				Project: targetParts[1],
				TagName: tagName,
			}
		}

		return &ghSource, nil
	}
}

// Function to get the github repo id from the repo information
func (g GitHubRepo) Id() string {
	return g.User + "/" + g.Project
}

// Function which gets the latest appimage from github release
func (g GitHubRepo) GetLatestRelease() (*Release, error) {
	client := github.NewClient(nil) // Client For Interacting with github api
	// Get all the releases from the target
	releases, _, err := client.Repositories.ListReleases(context.Background(), g.User, g.Project, nil)
	if err != nil {
		return nil, err
	}

	if g.TagName != "" {
		releaseWithTagName := getReleaseFromTagName(releases, g.TagName)

		if releaseWithTagName != nil {
			appimageFiles := getAppImageFilesFromRelease(releaseWithTagName)
			if len(appimageFiles) > 0 {
				return &Release{
					*releaseWithTagName.TagName,
					appimageFiles,
				}, nil
			}
		}	
	}

	// Filter out files which are not AppImage
	for _, release := range releases {
		if *release.Draft {
			continue
		}

		downloadLinks := getAppImageFilesFromRelease(release)
		if len(downloadLinks) > 0 {
			g.TagName = *release.TagName
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

func getAppImageFilesFromRelease(release *github.RepositoryRelease) ([]utils.BinaryUrl) {
	var downloadLinks []utils.BinaryUrl // Contains Download Links

	for _, asset := range release.Assets {
		if strings.HasSuffix(strings.ToLower(*asset.Name), ".appimage") {
			downloadLinks = append(downloadLinks, utils.BinaryUrl{
				FileName: *asset.Name,
				Url:      *asset.BrowserDownloadURL,
			})
		}
	}

	return downloadLinks
}

func getReleaseFromTagName(releases []*github.RepositoryRelease, tagName string) (*github.RepositoryRelease) {
	for _, release := range releases {
		if *release.Draft { continue }
		if tagName == *release.TagName {
			return release
		}
	}
	return nil
}
