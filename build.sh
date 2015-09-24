#! /bin/bash

set -e -u -o pipefail

cd "$(dirname "$BASH_SOURCE")"

# Always start fresh. Keeps things simple.
rm -rf build/

# Assemble the local GOPATH.
export GOPATH="$(pwd)/build"
mkdir -p "$GOPATH/src/github.com/keybase"
cp -r client "$GOPATH/src/github.com/keybase"

build_target="github.com/keybase/client/go/keybase"
echo Fetching dependencies with go get...
go get "$build_target"
echo Building kbstage...
go build -a -tags ./staging -o ./kbstage "$build_target"
