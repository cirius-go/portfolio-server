#!/usr/bin/env bash

set -e

BUILD_DIR="$PWD/cmd/infra/aws/.build/lambda"
mkdir -p "$BUILD_DIR"

first_clean() {
  rm -rf "$BUILD_DIR/*-worker.zip"

  mkdir -p "$BUILD_DIR/workers"
  rm -rf "$BUILD_DIR/workers/*"
}

build_worker() {
  local worker_name=$(basename "$1") # Extract worker name
  local worker_path="$1"
  local build_fn_dir="$BUILD_DIR/workers/$worker_name"
  local bootstrap_file="$build_fn_dir/bootstrap"
  local zip_file="$BUILD_DIR/$worker_name-worker.zip"

  echo "🔹 Building worker: $worker_name"

  mkdir -p "$build_fn_dir"

  # Set Lambda build environment
  export GOOS=linux
  export GOARCH=arm64
  export CGO_ENABLED=0

  # Build Go binary
  go build -trimpath -buildvcs=true -ldflags="-s -w" -o "$bootstrap_file" "$worker_path"

  # Package into ZIP
  cd "$build_fn_dir" || exit
  zip -0r "$zip_file" .
  echo "✅ Built: $zip_file"
}

main() {
  first_clean

  # Loop through all workers inside cmd/workers/
  for worker in "$PWD"/cmd/workers/*; do
    if [[ -d "$worker" ]]; then # Ensure it's a directory
      build_worker "$worker"
    fi
  done

  echo "🎉 All workers built successfully!"
}

main "$@"
