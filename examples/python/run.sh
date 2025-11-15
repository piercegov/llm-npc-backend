#!/bin/bash

# Run the example game client using uv
# This script automatically manages the virtual environment and dependencies

set -e

# Check if uv is installed
if ! command -v uv &> /dev/null; then
    echo "Error: uv is not installed"
    echo "Install it with: brew install uv"
    exit 1
fi

# Run the example using uv (it will automatically create venv and install deps)
echo "Running game client example..."
uv run game_client.py

