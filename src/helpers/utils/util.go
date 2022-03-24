package utils

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"os"
	"fmt"
	"bytes"
	"os/user"
	"net/url"
	"strings"
	"debug/elf"
	"path/filepath"
	"github.com/manifoldco/promptui"
)

type BinaryUrl struct {
	FileName string
	Url      string
}

type AppImageInfo struct {
	IsTerminalApp bool
	AppImageType int
}

// Get the user/repo from a github url
func GetUserRepoFromUrl(gitHubUrl string) (string, error) {
	urlParsed, err := url.ParseRequestURI(gitHubUrl)
	if err != nil {
		return "", err
	}

	if urlParsed.Host != "github.com" {
		return "", fmt.Errorf("invalid github url")
	}

	splitPaths := strings.Split(urlParsed.EscapedPath(), "/")

	if len(splitPaths) < 3 {
		return "", fmt.Errorf("invalid github url")
	}

	return splitPaths[1] + "/" + splitPaths[2], nil
}

// Get the Applications directory absolute path
func MakeApplicationsDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	applicationsPath := filepath.Join(usr.HomeDir, "Applications")
	err = os.MkdirAll(applicationsPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return applicationsPath, nil
}

// Get the file path of a file in applications folder
func MakeTargetFilePath(link *BinaryUrl) (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, link.FileName)
	return filePath, nil
}

// Make file path from a file in run-cache directory inside Applications directory
func MakeTempFilePath(link *BinaryUrl) (string, error) {
	applicationsPath, err := MakeTempAppDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, link.FileName)
	return filePath, nil
}

// Make folder run-cache inside Applications dir and return it's path
func MakeTempAppDirPath() (string, error) {
	TempApplicationDirPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	TempApplicationDirPath = filepath.Join(TempApplicationDirPath, "run-cache")
	err = os.MkdirAll(TempApplicationDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return TempApplicationDirPath, nil
}

// List appimages to select from
func PromptBinarySelection(downloadLinks []BinaryUrl) (result *BinaryUrl, err error) {
	if len(downloadLinks) == 1 {
		return &downloadLinks[0], nil
	}

	prompt := promptui.Select{
		Label: "Select an AppImage to install",
		Items: downloadLinks,
		Templates: &promptui.SelectTemplates{
			Label:    "   {{ .FileName }}",
			Active:   "\U00002713 {{ .FileName }}",
			Inactive: "   {{ .FileName }}",
			Selected: "\U00002713 {{ .FileName }}"},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &downloadLinks[i], nil
}

// read the update info embeded into the appimage file
func ReadUpdateInfo(appImagePath string) (string, error) {
	elfFile, err := elf.Open(appImagePath)
	if err != nil {
		panic("Unable to open target: \"" + appImagePath + "\"." + err.Error())
	}

	updInfo := elfFile.Section(".upd_info")
	if updInfo == nil {
		panic("Missing update section on target elf ")
	}
	sectionData, err := updInfo.Data()

	if err != nil {
		panic("Unable to parse update section: " + err.Error())
	}

	str_end := bytes.Index(sectionData, []byte("\000"))
	if str_end == -1 || str_end == 0 {
		return "", fmt.Errorf("No update information found in: " + appImagePath)
	}

	update_info := string(sectionData[:str_end])
	return update_info, nil
}
