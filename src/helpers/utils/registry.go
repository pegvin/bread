package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	updateUtils "github.com/DEVLOPRR/appimage-update/util"
)

type RegistryEntry struct {
	Repo       string
	FileSha1   string
	AppName    string
	AppVersion string
	FilePath   string
	UpdateInfo string
}

type Registry struct {
	Entries map[string]RegistryEntry
}

// Function to open a registry entry
func OpenRegistry() (registry *Registry, err error) {
	path, err := makeRegistryFilePath()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &Registry{Entries: map[string]RegistryEntry{}}, nil
	}

	err = json.Unmarshal(data, &registry)
	if err != nil {
		return
	}

	return
}

// Function to close a registry entry
func (registry *Registry) Close() error {
	path, err := makeRegistryFilePath()
	if err != nil {
		return err
	}

	blob, err := json.Marshal(registry)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, blob, 0666)
	if err != nil {
		return err
	}

	return nil
}

// Function to add a registry entry
func (registry *Registry) Add(entry RegistryEntry) error {
	registry.Entries[entry.FilePath] = entry
	return nil
}

// Function to remove a registry entry
func (registry *Registry) Remove(filePath string) {
	delete(registry.Entries, filePath)
}

// Function to update a registry entry
func (registry *Registry) Update() {
	applicationsDir, err := MakeApplicationsDirPath()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(applicationsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".AppImage") {
			filePath := filepath.Join(applicationsDir, f.Name())
			_, ok := registry.Entries[filePath]
			if !ok {
				entry := registry.createEntryFromFile(filePath)
				_ = registry.Add(entry)
			}
		}
	}

	// for filePath, _ := range registry.Entries {
	for filePath := range registry.Entries {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
			registry.Remove(filePath)
		}
	}
}

// func (registry *Registry) addFile(filePath string) {
// 	entry := registry.createEntryFromFile(filePath)
// 	_ = registry.Add(entry)
// }

// Function which creates a new entry in the registry from a file
func (registry *Registry) createEntryFromFile(filePath string) RegistryEntry {
	fileSha1, _ := GetFileSHA1(filePath)
	updateInfo, _ := updateUtils.ReadUpdateInfo(filePath)
	entry := RegistryEntry{
		Repo:       "",
		FileSha1:   fileSha1,
		AppName:    "",
		AppVersion: "",
		FilePath:   filePath,
		UpdateInfo: updateInfo,
	}
	return entry
}

// Function to lookup a entry in the registry
func (registry *Registry) Lookup(target string) (RegistryEntry, bool) {
	applicationsDir, _ := MakeApplicationsDirPath()
	possibleFullPath := filepath.Join(applicationsDir, target)

	for _, entry := range registry.Entries {
		if entry.FileSha1 == target || entry.FilePath == target ||
			entry.FilePath == possibleFullPath || entry.Repo == target {
			return entry, true
		}
	}

	if IsAppImageFile(target) {
		entry := registry.createEntryFromFile(target)
		_ = registry.Add(entry)

		return entry, true
	} else {
		if IsAppImageFile(possibleFullPath) {
			entry := registry.createEntryFromFile(target)
			_ = registry.Add(entry)

			return entry, true
		}
	}

	return RegistryEntry{}, false
}

// Function which makes the registry file path
func makeRegistryFilePath() (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, ".registry.json")
	return filePath, nil
}
