#!/usr/bin/env bash
set -ex

VERSION=$(git describe --always --tags --long)
PLATFORM=""

if [[ ${TRAVIS_OS_NAME} == 'linux' ]]; then
	PLATFORM="linux"
elif [[ ${TRAVIS_OS_NAME} == 'osx' ]]; then
	PLATFORM="darwin"
else
	PLATFORM="windows"
	exit 1
fi

env GO111MODULE=on make dnaNode-${PLATFORM}

set +x
echo "dnaNode-${PLATFORM}-amd64 |" $(md5sum dnaNode-${PLATFORM}-amd64|cut -d ' ' -f1)

