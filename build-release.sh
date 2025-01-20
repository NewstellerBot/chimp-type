#!/usr/bin/env bash

APP_NAME="myapp"
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
  # On Windows, add .exe extension
  if [ "$GOOS" = "windows" ]; then
    OUTPUT_NAME+='.exe'
  fi
  
  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o "dist/$OUTPUT_NAME" main.go
done

