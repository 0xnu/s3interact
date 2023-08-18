#!/bin/bash
# chmod +x package.sh && ./package.sh

# # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#                                                       #
#     created by: Finbarrs Oketunji <f@finbarrs.eu>     #
#     created on: 15/08/2023                            #
#                                                       #        
# # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# Define the OS and architecture combinations you want to build for
OS="linux darwin windows"
ARCH="amd64"

# Create the directory structure
for os in $OS; do
  for arch in $ARCH; do
    mkdir -p "build/${os}_${arch}"
  done
done

# Run gox
gox -os="$OS" -arch="$ARCH" -output="build/{{.OS}}_{{.Arch}}/s3interact-cli_{{.OS}}_{{.Arch}}"
