#!/bin/bash

# Build and run the log viewer (GUI or CLI mode)
echo "Building log viewer..."
go build ./cmd/logviewer/...

if [ $? -eq 0 ]; then
    if [ "$1" = "--cli" ] || [ "$1" = "-cli" ]; then
        echo "Starting log viewer in CLI mode..."
        ./logviewer --cli
    else
        echo "Starting log viewer GUI..."
        ./logviewer
    fi
else
    echo "Failed to build log viewer"
    exit 1
fi