#!/bin/bash

# This script tests GitHub Actions workflows locally using act
# https://github.com/nektos/act

set -e

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "Error: 'act' is not installed. Please install it first."
    echo "Installation instructions: https://github.com/nektos/act#installation"
    exit 1
fi

# Default workflow to test
WORKFLOW="release-action.yml"
DRYRUN=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -w|--workflow)
            WORKFLOW="$2"
            shift
            shift
            ;;
        -d|--dry-run)
            DRYRUN=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-w|--workflow WORKFLOW_FILE] [-d|--dry-run]"
            echo ""
            echo "Options:"
            echo "  -w, --workflow WORKFLOW_FILE  Specify the workflow file to test (default: release-action.yml)"
            echo "  -d, --dry-run                 Run in dry-run mode (no actual GitHub API calls)"
            echo "  -h, --help                    Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo "Testing GitHub Actions workflow: $WORKFLOW"

# Check if .secrets file exists
if [ ! -f .secrets ]; then
    echo "Error: .secrets file not found."
    echo "Please create a .secrets file with your GitHub token."
    echo "Example:"
    echo "GITHUB_TOKEN=your_github_token_here"
    exit 1
fi

# Source the .secrets file to get the GitHub token
source .secrets

# Validate GitHub token
if [ "$GITHUB_TOKEN" = "your_github_token_here" ] || [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: Invalid GitHub token in .secrets file."
    echo "Please update the GITHUB_TOKEN value in the .secrets file."
    exit 1
fi

echo "GitHub token found in .secrets file."

# Run the workflow with act
if [ "$DRYRUN" = true ]; then
    echo "Running in dry-run mode (no actual GitHub API calls)..."
    act workflow_dispatch -W .github/workflows/$WORKFLOW --secret-file .secrets --dryrun
else
    echo "Running workflow with act..."
    act workflow_dispatch -W .github/workflows/$WORKFLOW --secret-file .secrets --verbose
fi

echo "Workflow test completed!" 