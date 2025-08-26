#!/bin/bash

# Script to create GitHub release using GitHub CLI
# Make sure you're authenticated: gh auth login

set -e

VERSION="v0.1.1"
RELEASE_FILE="p3ipam-$VERSION-darwin-amd64.tar.gz"
RELEASE_NOTES="RELEASE_NOTES.md"

echo "🚀 Creating GitHub Release for p3ipam $VERSION"
echo "=============================================="

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed"
    echo "Install with: brew install gh"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub"
    echo "Run: gh auth login"
    exit 1
fi

# Check if release file exists
if [ ! -f "$RELEASE_FILE" ]; then
    echo "Error: Release file $RELEASE_FILE not found"
    exit 1
fi

# Check if release notes exist
if [ ! -f "$RELEASE_NOTES" ]; then
    echo "Error: Release notes $RELEASE_NOTES not found"
    exit 1
fi

echo "📦 Creating release on GitHub..."
echo "Version: $VERSION"
echo "File: $RELEASE_FILE"
echo "Notes: $RELEASE_NOTES"

# Create the release
gh release create "$VERSION" \
    --title "p3ipam $VERSION - Binary-Only Release" \
    --notes-file "$RELEASE_NOTES" \
    "$RELEASE_FILE"

echo ""
echo "✅ Release created successfully!"
echo "🎉 Check it out at: https://github.com/palmar/p3ipam/releases"
