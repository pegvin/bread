#!/bin/bash

set -e

COMPILER="go"
BINARY="bread"
DIST="build"
ENTRY_FILE="src/main.go"
OUTPUT="$DIST/$BINARY"
VERSION="0.2.2"

if [[ $1 = '' ]]; then
	echo "Compiling '$ENTRY_FILE' into '$DIST'"
	${COMPILER} build -o ${OUTPUT} -v ${ENTRY_FILE}
	echo "Compiled Successfully into '$OUTPUT'"
elif [[ $1 = 'appimage' ]]; then
	echo "Building AppImage"
	BREAD_VERSION=$VERSION appimage-builder --skip-test --recipe=AppImage-Builder.yml
elif [[ $1 = 'get-deps' ]]; then
	echo "Getting Depedencies"
	${COMPILER} get -v -t -d ./...
	${COMPILER} mod tidy
elif [[ $1 = 'clean' ]]; then
	rm -rfv $DIST
	rm -rfv appimage-builder-cache
	rm -rfv AppDir
	rm -rfv bread-*.AppImage*
else
	echo "Build Script '$1' Not Found!"
fi
