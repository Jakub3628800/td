#!/bin/bash

# This script tests the release workflow locally
# It creates a mock release without actually pushing to GitHub

set -e

# Check if required tools are installed
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' is not installed. Please install it first."
    echo "Installation: sudo apt-get install jq"
    exit 1
fi

# Extract current version
CURRENT_VERSION=$(grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
echo "Current version: $CURRENT_VERSION"

# Simulate version bump (without actually changing the file)
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="${MAJOR}.${MINOR}.${NEW_PATCH}"
echo "Simulated new version: $NEW_VERSION"

# Build release binaries
echo "Building release binaries..."
make build-release

# Create a mock release
echo "Creating mock release for v$CURRENT_VERSION..."
echo "This would create a GitHub release with the following files:"
ls -la bin/td-$CURRENT_VERSION-*

echo ""
echo "To create an actual release, you would need to:"
echo "1. Go to GitHub repository"
echo "2. Click on 'Actions' tab"
echo "3. Select 'Create Release' workflow"
echo "4. Click 'Run workflow'"
echo "5. Select version type and whether it's a prerelease"
echo "6. Click 'Run workflow'"

echo ""
echo "For local testing with act, make sure your .secrets file contains a valid GitHub token with 'repo' permissions."
echo "Then run: ./scripts/test-github-actions.sh" 