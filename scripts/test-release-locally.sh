#!/bin/bash

# This script tests the release workflow locally
# It creates a mock release without actually pushing to GitHub

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print colored message
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Check if required tools are installed
if ! command -v jq &> /dev/null; then
    print_message "$YELLOW" "Warning: 'jq' is not installed. Some features may be limited."
    print_message "$YELLOW" "Installation: sudo apt-get install jq"
fi

# Parse command line arguments
VERSION_TYPE="patch"
SKIP_BUILD=false

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -t|--type)
            VERSION_TYPE="$2"
            shift
            shift
            ;;
        -s|--skip-build)
            SKIP_BUILD=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-t|--type VERSION_TYPE] [-s|--skip-build]"
            echo ""
            echo "Options:"
            echo "  -t, --type VERSION_TYPE  Specify the version type to bump (major, minor, patch)"
            echo "  -s, --skip-build         Skip building the binaries"
            echo "  -h, --help               Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Validate version type
if [[ "$VERSION_TYPE" != "major" && "$VERSION_TYPE" != "minor" && "$VERSION_TYPE" != "patch" ]]; then
    print_message "$RED" "Error: Invalid version type. Use 'major', 'minor', or 'patch'."
    exit 1
fi

# Extract current version
CURRENT_VERSION=$(grep -o '"[0-9]\+\.[0-9]\+\.[0-9]\+"' main.go | tr -d '"')
print_message "$GREEN" "Current version: $CURRENT_VERSION"

# Simulate version bump (without actually changing the file)
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
if [[ "$VERSION_TYPE" == "major" ]]; then
    NEW_MAJOR=$((MAJOR + 1))
    NEW_VERSION="${NEW_MAJOR}.0.0"
elif [[ "$VERSION_TYPE" == "minor" ]]; then
    NEW_MINOR=$((MINOR + 1))
    NEW_VERSION="${MAJOR}.${NEW_MINOR}.0"
else
    NEW_PATCH=$((PATCH + 1))
    NEW_VERSION="${MAJOR}.${MINOR}.${NEW_PATCH}"
fi
print_message "$GREEN" "Simulated new version: $NEW_VERSION (${VERSION_TYPE} bump)"

# Build release binaries
if [ "$SKIP_BUILD" = false ]; then
    print_message "$GREEN" "Building release binaries..."
    make build-release || {
        print_message "$RED" "Error: Failed to build release binaries."
        exit 1
    }
else
    print_message "$YELLOW" "Skipping binary build..."
fi

# Create a mock release
print_message "$GREEN" "Creating mock release for v$CURRENT_VERSION..."
print_message "$GREEN" "This would create a GitHub release with the following files:"
ls -la bin/td-$CURRENT_VERSION-* 2>/dev/null || print_message "$YELLOW" "No matching binaries found. Run without --skip-build to create them."

echo ""
print_message "$GREEN" "To create an actual release, you would need to:"
echo "1. Go to GitHub repository"
echo "2. Click on 'Actions' tab"
echo "3. Select 'Create Release' workflow"
echo "4. Click 'Run workflow'"
echo "5. Select version type: $VERSION_TYPE"
echo "6. Click 'Run workflow'"

echo ""
print_message "$GREEN" "For local testing with act:"
echo "1. Make sure your .secrets file contains a valid GitHub token with 'repo' permissions."
echo "2. Run: ./scripts/test-github-actions.sh -w local-release-test.yml"
echo "3. For a dry run (no actual API calls): ./scripts/test-github-actions.sh -w local-release-test.yml --dry-run" 