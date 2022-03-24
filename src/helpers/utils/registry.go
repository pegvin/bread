package utils

import (
	"os"
	"log"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	updateUtils "github.com/DEVLOPRR/appimage-update/util"
)

type RegistryEntry struct {
	Repo          string
	FileSha1      string
	AppName       string
	AppVersion    string
	FilePath      string
	UpdateInfo    string
	IsTerminalApp bool
	AppImageType  int
}

type Registry struct {
	Entries map[string]RegistryEntry
}

// Function to open a registry entry
func OpenRegistry() (registry *Registry, err error) {
	path, err := makeRegistryFilePath() // Get the full path to .registry.json
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(path) // Read file
	if err != nil {
		return &Registry{Entries: map[string]RegistryEntry{}}, nil // If some error occured return a new empty registry
	}

	err = json.Unmarshal(data, &registry) // Parse JSON data into the struct
	if err != nil {
		return
	}

	return
}

// Function to close a registry entry
func (registry *Registry) Close() error {
	path, err := makeRegistryFilePath() // Get full path to .registry.json
	if err != nil {
		return err
	}

	blob, err := json.Marshal(registry) // Convert registry struct into a blob
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, blob, 0666) // write the blob file with 666 permissions
	if err != nil {
		return err
	}

	return nil
}

// Add a entry to registry
func (registry *Registry) Add(entry RegistryEntry) error {
	registry.Entries[entry.FilePath] = entry
	return nil
}

// Remove a entry from registry
func (registry *Registry) Remove(filePath string) {
	delete(registry.Entries, filePath)
}

// Update registry entry
func (registry *Registry) Update() {
	applicationsDir, err := MakeApplicationsDirPath() // Applications folder full path
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(applicationsDir) // Read all the files in the folder
	if err != nil {
		log.Fatal(err)
	}

	// Filter out all the appimage files and put them into registry
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".appimage") {
			filePath := filepath.Join(applicationsDir, f.Name())
			_, ok := registry.Entries[filePath]
			if !ok {
				entry := registry.createEntryFromFile(filePath)
				_ = registry.Add(entry)
			}
		}
	}

	// Remove deleted/non-existent files from registry
	for filePath := range registry.Entries {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
			registry.Remove(filePath)
		}
	}
}

// Create a new entry in the registry from a appimage file
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

// Lookup a entry in the registry
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

// makes the registry file path
func makeRegistryFilePath() (string, error) {
	applicationsPath, err := MakeApplicationsDirPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(applicationsPath, ".registry.json")
	return filePath, nil
}
