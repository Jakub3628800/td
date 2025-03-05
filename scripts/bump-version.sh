#!/bin/bash

# This script bumps the version in main.go
# Usage: ./scripts/bump-version.sh [major|minor|patch]

set -e

# Check if the version type is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 [major|minor|patch]"
    exit 1
fi

VERSION_TYPE=$1

# Check if the version type is valid
if [ "$VERSION_TYPE" != "major" ] && [ "$VERSION_TYPE" != "minor" ] && [ "$VERSION_TYPE" != "patch" ]; then
    echo "Invalid version type. Use 'major', 'minor', or 'patch'."
    exit 1
fi

# Get the current version from main.go
CURRENT_VERSION=$(grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')

# Split the version into major, minor, and patch
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Bump the version based on the version type
if [ "$VERSION_TYPE" == "major" ]; then
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
elif [ "$VERSION_TYPE" == "minor" ]; then
    MINOR=$((MINOR + 1))
    PATCH=0
else
    PATCH=$((PATCH + 1))
fi

# Create the new version
NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"

# Update the version in main.go
sed -i "s/Version = \"[0-9]\+\.[0-9]\+\.[0-9]\+\"/Version = \"${NEW_VERSION}\"/" main.go

echo "Version bumped from $CURRENT_VERSION to $NEW_VERSION"

# Commit the changes
git add main.go
git commit -m "Bump version to $NEW_VERSION"

echo "Changes committed. You can now push the changes with 'git push'." 