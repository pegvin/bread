#!/bin/bash

set -e

COMPILER="go"
BINARY="bread"
DIST="build"
ENTRY_FILE="src/main.go"
OUTPUT="$DIST/$BINARY"
VERSION="$(cat src/main.go | grep '"VERSION":' | grep -o '[0-9 .]*' | xargs)" # Simple Hack To Get The Version Number From main.go file

if [[ $1 = '' || $1 = '--prod' ]]; then
	echo "Compiling '$ENTRY_FILE' into '$DIST'"
	if [[ $1 = '--prod' ]]; then
		${COMPILER} build -ldflags "-s -w" -o ${OUTPUT} -v ${ENTRY_FILE}
	else
		${COMPILER} build -o ${OUTPUT} -v ${ENTRY_FILE}
	fi
	echo "Compiled Successfully into '$OUTPUT'"
elif [[ $1 = 'appimage' ]]; then
	echo "Building AppImage"
	BREAD_VERSION=$VERSION appimage-builder --skip-test --recipe=AppImage-Builder.yml
elif [[ $1 = 'get-deps' ]]; then
	echo "Getting Depedencies"
	${COMPILER} mod tidy
elif [[ $1 = 'clean' ]]; then
	rm -rfv $DIST
	rm -rfv appimage-builder-cache
	rm -rfv AppDir
	rm -rfv bread-*.AppImage*
else
	echo "Build Script '$1' Not Found!"
fi
