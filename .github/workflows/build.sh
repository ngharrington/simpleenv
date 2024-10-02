#!/bin/bash

function is_semver() {
  local version_string="$1"
  local semver_pattern='^testv([0-9]+)\.([0-9]+)\.([0-9]+)$'

  if [[ $version_string =~ $semver_pattern ]]; then
    return 0
  else
    return 1
  fi
}

# Parse the version from the command line arg
VERSION=$1
if ! is_semver $VERSION; then
  echo "Invalid version: $VERSION"
  exit 1
fi


arch="amd64"

GOARCH=$arch go build -o ./ ./...

release="simpleenv-${VERSION}-${arch}"

mkdir -p $release

cp ./simpleenv $release/

tarball="${release}.tar.gz"

tar -czf $tarball $release

if [ $? -ne 0 ]; then
	  echo "Failed to create $tarball"
	  exit 1
fi

