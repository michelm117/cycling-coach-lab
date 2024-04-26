#!/bin/bash

# Define the source directories
SOURCE_DIRS=("services" "utils" "db")

# Define the destination directory
DEST_DIR="mocks"

# Iterate over each source directory
for SOURCE_DIR in "${SOURCE_DIRS[@]}"; do
    # Find all source files excluding those ending with _test.go
    SOURCE_FILES=$(find "$SOURCE_DIR" -type f -name '*.go' ! -name '*_test.go')

    # Iterate over each source file
    for SOURCE_FILE in $SOURCE_FILES; do
        # Extract the file name without the directory path
        FILENAME=$(basename "$SOURCE_FILE")

        # Generate the destination file path by replacing the source directory with the destination directory
        DEST_FILE="${DEST_DIR}/${FILENAME}"

        echo "Generating mock for $SOURCE_FILE to $DEST_FILE"
        mockgen -source "$SOURCE_FILE" -destination "$DEST_FILE" -package mocks
    done
done

