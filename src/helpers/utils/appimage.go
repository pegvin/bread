package utils

import (
	"os"
	"bytes"
)

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
