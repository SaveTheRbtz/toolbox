#!/usr/bin/env bash

# Prerequisite: set absolute path to project root to $PROJECT_ROOT

BUILD_PATH=/build
DIST_PATH=/dist

BUILD_MAJOR_VERSION=$(cat $PROJECT_ROOT/version)
BUILD_HASH=$(cd $PROJECT_ROOT && git rev-parse HEAD)

if [ ! -d $BUILD_PATH ]; then
  mkdir -p $BUILD_PATH
fi
if [ "$TOOLBOX_BUILD_ID"x = ""x ]; then
  TOOLBOX_BUILD_ID=0.0
fi
BUILD_VERSION=$BUILD_MAJOR_VERSION.$TOOLBOX_BUILD_ID

echo --------------------
echo BUILD: Start building version: $BUILD_VERSION

echo --------------------
echo BUILD: Testing..

cd $PROJECT_ROOT
go test $(glide novendor)


echo --------------------
echo BUILD: Building tool

X_APP_NAME="-X github.com/watermint/toolbox/infra.AppName=toolbox"
X_APP_VERSION="-X github.com/watermint/toolbox/infra.AppVersion=$BUILD_VERSION"
X_APP_HASH="-X github.com/watermint/toolbox/infra.AppHash=$BUILD_HASH"
X_APP_CREDENTIALS=""
if [ -e $PROJECT_ROOT/credentials.secret ]; then
X_APP_CREDENTIALS=$(cat $PROJECT_ROOT/credentials.secret | xargs)
fi
if [ "$TOOLBOX_APP_CREDENTIALS"x != ""x ]; then
X_APP_CREDENTIALS="$TOOLBOX_APP_CREDENTIALS"
fi
LD_FLAGS="$X_APP_NAME $X_APP_VERSION $X_APP_HASH $X_APP_CREDENTIALS"

cd $BUILD_PATH
GOOS=windows GOARCH=386   go build --ldflags "$LD_FLAGS" -o tbx-$BUILD_VERSION-win.exe github.com/watermint/toolbox
GOOS=linux   GOARCH=386   go build --ldflags "$LD_FLAGS" -o tbx-$BUILD_VERSION-linux   github.com/watermint/toolbox
GOOS=darwin  GOARCH=amd64 go build --ldflags "$LD_FLAGS" -o tbx-$BUILD_VERSION-macos   github.com/watermint/toolbox


echo --------------------
echo BUILD: Packaging
zip -9 -r $DIST_PATH/tbx-$BUILD_VERSION.zip tbx-*