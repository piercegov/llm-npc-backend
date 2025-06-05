#!/bin/bash

# Build and run the log viewer GUI
echo "Building log viewer..."
go build ./cmd/logviewer/...

if [ $? -eq 0 ]; then
    echo "Starting log viewer GUI..."
    ./logviewer
else
    echo "Failed to build log viewer"
    exit 1
fi