package utils

import (
	"os"
	"bytes"
	"github.com/DEVLOPRR/libappimage-go"
)

// Get AppImage information: isTerminalApp, AppImageType
func GetAppImageInfo(targetFilePath string, debug bool) (*AppImageInfo, error) {
	libAppImage, err := libappimagego.NewLibAppImageBindings() // Load the `libappimage` Library For Integration
	if err != nil {
		return nil, err
	}

	return &AppImageInfo{
		IsTerminalApp: libAppImage.IsTerminalApp(targetFilePath),
		AppImageType: libAppImage.GetType(targetFilePath, debug),
	}, nil
}

// Remove the application desktop integration
func RemoveDesktopIntegration(filePath string, debug bool) (error) {
	libAppImage, err := libappimagego.NewLibAppImageBindings()
	if err != nil {
		return err
	}

	if libAppImage.ShallNotBeIntegrated(filePath) {
		return nil
	}

	err = libAppImage.Unregister(filePath, debug)
	return err
}

// Integrate The AppImage To Desktop.
func CreateDesktopIntegration(targetFilePath string, debug bool) (error) {
	libAppImage, err := libappimagego.NewLibAppImageBindings() // Load the `libappimage` Library For Integration
	if err != nil {
		return err
	}

	err = libAppImage.Register(targetFilePath, debug) // Register The File
	if err != nil {
		return err
	} else {
		return nil
	}
}

// check if a file is appimage
func IsAppImageFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}

	return isAppImageType1File(file) || isAppImageType2File(file)
}

// Check if appimage is type 2
func isAppImageType2File(file *os.File) bool {
	return isElfFile(file) && fileHasBytesAt(file, []byte{0x41, 0x49, 0x02}, 8)
}

// Check if appimage is type 1
func isAppImageType1File(file *os.File) bool {
	return isISO9660(file) || fileHasBytesAt(file, []byte{0x41, 0x49, 0x01}, 8)
}

// Check if a file is Elf file
func isElfFile(file *os.File) bool {
	return fileHasBytesAt(file, []byte{0x7f, 0x45, 0x4c, 0x46}, 0)
}

// Check if the file is a ISO 9660 Standard File
func isISO9660(file *os.File) bool {
	for _, offset := range []int64{32769, 34817, 36865} {
		if fileHasBytesAt(file, []byte{'C', 'D', '0', '0', '1'}, offset) {
			return true
		}
	}

	return false
}

// check if a file has bytes at particular position
func fileHasBytesAt(file *os.File, expectedBytes []byte, offset int64) bool {
	readBytes := make([]byte, len(expectedBytes))
	_, _ = file.Seek(offset, 0)
	_, _ = file.Read(readBytes)

	return bytes.Equal(readBytes, expectedBytes)
}
