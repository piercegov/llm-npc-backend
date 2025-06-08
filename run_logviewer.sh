#!/bin/bash

# Build and run the log viewer (GUI or CLI mode)
echo "Building log viewer..."
go build ./cmd/logviewer/...

if [ $? -eq 0 ]; then
    # Parse arguments to determine mode and backend options
    CLI_MODE=false
    HTTP_MODE=false
    
    for arg in "$@"; do
        case $arg in
            --cli|-cli)
                CLI_MODE=true
                ;;
            --http)
                HTTP_MODE=true
                ;;
        esac
    done
    
    # Build logviewer arguments
    LOGVIEWER_ARGS=""
    if [ "$CLI_MODE" = true ]; then
        LOGVIEWER_ARGS="$LOGVIEWER_ARGS --cli"
    fi
    if [ "$HTTP_MODE" = true ]; then
        LOGVIEWER_ARGS="$LOGVIEWER_ARGS --http"
    fi
    
    if [ "$CLI_MODE" = true ]; then
        echo "Starting log viewer in CLI mode..."
        if [ "$HTTP_MODE" = true ]; then
            echo "Backend will run in HTTP mode"
        else
            echo "Backend will run in Unix socket mode"
        fi
    else
        echo "Starting log viewer GUI..."
        if [ "$HTTP_MODE" = true ]; then
            echo "Backend will run in HTTP mode"
        else
            echo "Backend will run in Unix socket mode"
        fi
    fi
    
    ./logviewer $LOGVIEWER_ARGS
else
    echo "Failed to build log viewer"
    exit 1
fi