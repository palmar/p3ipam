#!/bin/bash

# p3ipam Linux Release Creation Script
# This script creates a clean Linux release without macOS metadata

set -e

VERSION="v0.1.1"
RELEASE_NAME="p3ipam-$VERSION-linux-amd64"

echo "🐧 Creating p3ipam $VERSION Linux Release"
echo "========================================="

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "Error: Please run this script from the p3ipam root directory"
    exit 1
fi

echo "🔨 Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o p3ipam-linux-amd64 main.go

echo "📦 Preparing clean Linux release..."

# Create a clean release directory
mkdir -p "linux-release"
cp p3ipam-linux-amd64 "linux-release/p3ipam"
cp schema.sql "linux-release/"

# Ensure clean permissions
chmod +x "linux-release/p3ipam"
chmod 644 "linux-release/schema.sql"

# Create release archive
cd "linux-release"
tar -czf "../$RELEASE_NAME.tar.gz" .
cd ..

# Clean up temporary directory
rm -rf "linux-release"
rm p3ipam-linux-amd64

# Create checksum
shasum -a 256 "$RELEASE_NAME.tar.gz" > "$RELEASE_NAME.tar.gz.sha256"

echo "✅ Linux release created: $RELEASE_NAME.tar.gz"
echo "✅ Checksum created: $RELEASE_NAME.tar.gz.sha256"

# Show archive contents
echo ""
echo "📁 Archive contents:"
tar -tzf "$RELEASE_NAME.tar.gz"

echo ""
echo "🎯 Ready to upload to GitHub release!"
echo "Run: gh release upload v0.1.1 $RELEASE_NAME.tar.gz"
