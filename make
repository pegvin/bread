#!/bin/bash

set -eo pipefail

COMPILER="go" # Compiler To Use For Building
BINARY="bread" # Output Binary Name
DIST="build" # Output Directory Name

# Simple Hack To Get The Version Number From main.go file
VERSION="0.7.2"
ENTRY_FILE="src/main.go" # Main Entry File To Compile
OUTPUT="$DIST/$BINARY" # Output Path Of Built Binary
COMPRESSED_OUTPUT="$OUTPUT-$VERSION-x86_64" # Output path of the compressed binary
COMMIT_HASH=$(git log --pretty=format:'%h' -n 1)

if [[ $1 = '' || $1 = '--prod' ]]; then
	echo "Compiling '$ENTRY_FILE' into '$DIST'"
	if [[ $1 = '--prod' ]]; then
		# When building for production use some ldflags and upx to reduce the binary size
		${COMPILER} build -ldflags "-s -w -X main.VERSION=${VERSION}" -o ${OUTPUT} -v ${ENTRY_FILE}
		upx -9 -o ${COMPRESSED_OUTPUT} ${OUTPUT}
	else
		${COMPILER} build -ldflags "-X main.VERSION=${VERSION}" -o ${OUTPUT} -v ${ENTRY_FILE}
	fi
	echo "Compiled Successfully into '$OUTPUT'"
elif [[ $1 = 'appimage' ]]; then
	echo "Building AppImage"
	# Set the bread version to a env variable and call appimage-builder to make the appimage
	BREAD_VERSION=$VERSION appimage-builder --skip-test --recipe=AppImage-Builder.yml
elif [[ $1 = 'get-deps' ]]; then
	echo "Getting Depedencies"
	${COMPILER} mod tidy
	go get -t -u ./...
	echo 'Done!'
elif [[ $1 = 'clean' ]]; then
	rm -rfv $DIST
	rm -rfv appimage-builder-cache
	rm -rfv AppDir
	rm -rfv bread-*.AppImage*
	echo 'Done!'
elif [[ $1 = 'reg' ]]; then
	# Simple hacky way to pretty print `.registry.json` in user's home dir
	node -e 'console.log(JSON.parse(require("fs").readFileSync(`${require("os").homedir()}/Applications/.registry.json`)))'
else
	echo "Build Script '$1' Not Found!"
fi
