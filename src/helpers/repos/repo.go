package repos

import (
	"bread/src/helpers/utils"
)

type Release struct {
	Tag   string
	Files []utils.BinaryUrl
}

type Application interface {
	Id() string
	GetLatestRelease() (*Release, error)
	Download(binaryUrl *utils.BinaryUrl, targetPath string) error
	FallBackUpdateInfo() string
}

// Function Which Parses String And Returns A Repo Object And Error (nil if not any)
func ParseTarget(target string) (Application, error) {
	// Parse The Repo As A GitHub Repo, And if there is no error return repo
	repo, err := NewGitHubRepo(target)
	if err == nil {
		return repo, nil
	}

	return nil, InvalidTargetFormat
}
