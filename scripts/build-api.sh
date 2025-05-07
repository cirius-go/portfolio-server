#!/usr/bin/env bash

set -e

BUILD_DIR="$PWD/cmd/infra/aws/.build/lambda"
BUILD_FN_DIR="$BUILD_DIR/api"
CMD_DIR="$PWD/cmd/api"
mkdir -p "$BUILD_DIR"

prepare_build_directory() {
  mkdir -p "$BUILD_FN_DIR"
  rm -rf "$BUILD_FN_DIR"/*
}

set_build_args() {
  buildArgs=()
  if [[ "$TARGET" == "lambda" ]]; then
    buildArgs+=("lambda" "lambda.norpc")
  fi
  echo "${buildArgs[@]}"
}

build_project() {
  local buildTags="$1"
  local bootstrapFile="$BUILD_FN_DIR/bootstrap"

  if [[ "$TARGET" == "lambda" ]]; then
    export GOOS=linux
    export GOARCH=arm64
    export CGO_ENABLED=0
  fi

  if [[ -n "$buildTags" ]]; then
    go build -trimpath -buildvcs=true -tags="$buildTags" -ldflags="-s -w" -o "$bootstrapFile" "$CMD_DIR"
  else
    go build -trimpath -buildvcs=true -ldflags="-s -w" -o "$bootstrapFile" "$CMD_DIR"
  fi
}

prepare_assets() {
  mkdir -p "$BUILD_FN_DIR/docs/swagger/"
  cp "$PWD"/docs/swagger/* "$BUILD_FN_DIR/docs/swagger"
  cp -r "$PWD/assets" "$BUILD_FN_DIR"
}

zip_package() {
  cd "$BUILD_FN_DIR" || exit
  zip -0r "../api.zip" .
}

main() {
  prepare_build_directory
  buildTags=$(set_build_args)
  build_project "$buildTags"
  prepare_assets
  zip_package
}

main "$@"
