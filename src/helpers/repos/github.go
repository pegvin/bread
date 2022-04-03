package repos

import (
	"context"
	"strings"

	"bread/src/helpers/utils"
	"github.com/google/go-github/v31/github"
)

// Struct containing GitHub Repo Details
type GitHubRepo struct {
	User     string
	Project  string
	Release  string
	File     string
	TagName  string
	UserRepo string
}

// Parses string to a github repo information, and returns a object and error (if any)
func NewGitHubRepo(target string, tagName string) (appInfo Application, err error) {
	appInfo = &GitHubRepo{}
	ghSource := GitHubRepo{}

	// parse the target as a github url and get the user/repo from it
	userRepo, err := utils.GetUserRepoFromUrl(target)
	if err == nil { // If successfull return the information
		userRepoSplitted := strings.Split(userRepo, "/")
		ghSource = GitHubRepo{
			User:    userRepoSplitted[0],
			Project: userRepoSplitted[1],
			TagName: tagName,
			UserRepo: userRepoSplitted[0] + "/" + userRepoSplitted[1],
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
				UserRepo: targetParts[0] + "/" + targetParts[0],
			}
		} else {
			ghSource = GitHubRepo{
				User:    targetParts[0],
				Project: targetParts[1],
				TagName: tagName,
				UserRepo: targetParts[0] + "/" + targetParts[1],
			}
		}

		return &ghSource, nil
	}
}

// Get the github user/repo from the repo information
func (g GitHubRepo) Id() string {
	return g.UserRepo
}

// Gets the latest/specified tagged release from github
func (g GitHubRepo) GetLatestRelease(NoPreRelease bool) (*Release, error) {
	client := github.NewClient(nil) // Client For Interacting with github api
	client.RateLimits(context.Background())
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
		if *release.Prerelease && NoPreRelease {
			continue
		}

		downloadLinks := getAppImageFilesFromRelease(release)
		if len(downloadLinks) > 0 {
			return &Release{
				*release.TagName,
				downloadLinks,
			}, nil
		}
	}

	return nil, NoAppImageBinariesFound
}

// Download appimage from remote
func (g GitHubRepo) Download(binaryUrl *utils.BinaryUrl, targetPath string) (err error) {
	err = utils.DownloadFile(binaryUrl.Url, targetPath, 0755, "Downloading")
	return err
}

// Generate a fallback update information for a appimage
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

// Gets All The AppImage Files from a github release
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

// Gets Release From A Particular Tag Name
func getReleaseFromTagName(releases []*github.RepositoryRelease, tagName string) (*github.RepositoryRelease) {
	for _, release := range releases {
		if *release.Draft { continue }
		if tagName == *release.TagName {
			return release
		}
	}
	return nil
}
