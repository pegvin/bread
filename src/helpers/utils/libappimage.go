package utils

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"github.com/rainycape/dl"
)

type libAppImageBind struct {
	lib *dl.DL

	appimage_shall_not_be_integrated func(path *C.char) int
	appimage_register_in_system      func(path *C.char, verbose int) int
	appimage_unregister_in_system    func(path *C.char, verbose int) int
}

type LibAppImage interface {
	Register(filePath string) error
	Unregister(filePath string) error
	ShallNotBeIntegrated(filePath string) bool
	Close()
}

// Function which loads up libappimage from the system, libappimage comes packed with the imageHub AppImage.
func loadLibAppImage() (*dl.DL, error) {
	// libappimage versions from latest to oldest, so that we can load the latest version
	sharedLibList := [17]string{
		".1.0.4", ".1.0.3", ".1.0.1", ".1.0.2",
		".1.0", ".0.1.9", ".0.1.8", ".0.1.7", ".0.1.6", ".0.1.5",
		".0.1.4", ".0.1.3", ".0.1.2", ".0.1.1", ".0.1.0", ".0", "",
	}

	for index := range sharedLibList {
		lib, err := dl.Open("libappimage.so" + sharedLibList[index], 0)
		if err == nil {
			return lib, nil
		}
	}

	return nil, fmt.Errorf("libappimage not found, desktop integration is disabled")
}

// Makes a new binding with libappimage, and returns a object with functions to register, unregister and other functions
func NewLibAppImageBindings() (LibAppImage, error) {
	bindings := libAppImageBind{}
	var err error
	bindings.lib, err = loadLibAppImage()

	if err != nil {
		return nil, err
	}

	err = bindings.lib.Sym("appimage_shall_not_be_integrated", &bindings.appimage_shall_not_be_integrated)
	if err != nil {
		return nil, err
	}

	err = bindings.lib.Sym("appimage_unregister_in_system", &bindings.appimage_unregister_in_system)
	if err != nil {
		return nil, err
	}

	err = bindings.lib.Sym("appimage_register_in_system", &bindings.appimage_register_in_system)
	if err != nil {
		return nil, err
	}

	return &bindings, nil
}

// Function to register a appimage
func (bind *libAppImageBind) Register(filePath string) error {
	if bind.appimage_register_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("unregister failed")
	}

	return nil
}

// Function to unregister a appimage
func (bind *libAppImageBind) Unregister(filePath string) error {
	if bind.appimage_unregister_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("unregister failed")
	}

	return nil
}

// Function which returns if a appimage should be integrated or not
func (bind *libAppImageBind) ShallNotBeIntegrated(filePath string) bool {
	return bind.appimage_shall_not_be_integrated(C.CString(filePath)) != 0
}

// Function to close the binding
func (bind *libAppImageBind) Close() {
	_ = bind.lib.Close()
}
