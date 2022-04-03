package repos

import (
	"strings"
	"bread/src/helpers/utils"
)

type Release struct {
	Tag   string
	Files []utils.BinaryUrl
}

type Application interface {
	Id() string
	GetLatestRelease(NoPreRelease bool) (*Release, error)
	Download(binaryUrl *utils.BinaryUrl, targetPath string) error
	FallBackUpdateInfo() string
}

// Parse String And Returns A Repo Object And Error (nil if not any)
func ParseTarget(target string, tagName string) (Application, error) {
	target = strings.ToLower(target)
	// Parse The Repo As A GitHub Repo, And if there is no error return repo
	repo, err := NewGitHubRepo(target, tagName)
	if err == nil {
		return repo, nil
	}

	return nil, InvalidTargetFormat
}
