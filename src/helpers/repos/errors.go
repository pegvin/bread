package repos

import "errors"

var InvalidTargetFormat = errors.New("invalid target format")

var NoAppImageBinariesFound = errors.New("no appImage found in releases")
