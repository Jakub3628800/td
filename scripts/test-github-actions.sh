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

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -w|--workflow)
            WORKFLOW="$2"
            shift
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-w|--workflow WORKFLOW_FILE]"
            echo ""
            echo "Options:"
            echo "  -w, --workflow WORKFLOW_FILE  Specify the workflow file to test (default: release-action.yml)"
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

# Run the workflow with act
# Note: This will use the .actrc file if it exists
act workflow_dispatch -W .github/workflows/$WORKFLOW --secret-file .secrets --verbose

echo "Workflow test completed!" 