#!/bin/bash

# p3ipam Release Creation Script
# This script helps create a GitHub release

set -e

VERSION="v0.1.1"
RELEASE_DIR="releases/$VERSION"

echo "🚀 Creating p3ipam $VERSION Release"
echo "=================================="

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "Error: Please run this script from the p3ipam root directory"
    exit 1
fi

# Check if release directory exists
if [ ! -d "$RELEASE_DIR" ]; then
    echo "Error: Release directory $RELEASE_DIR not found"
    exit 1
fi

echo "📦 Preparing release files..."

# Create a clean release directory with just the binary and essential files
mkdir -p "release-binary"
cp "$RELEASE_DIR/p3ipam" "release-binary/"
cp "$RELEASE_DIR/schema.sql" "release-binary/"

# Create release archive with just the binary and essential files
cd "release-binary"
tar -czf "../p3ipam-$VERSION-darwin-amd64.tar.gz" .
cd ..

# Clean up temporary directory
rm -rf "release-binary"

# Create checksum
shasum -a 256 "p3ipam-$VERSION-darwin-amd64.tar.gz" > "p3ipam-$VERSION-darwin-amd64.tar.gz.sha256"

echo "✅ Release archive created: p3ipam-$VERSION-darwin-amd64.tar.gz"
echo "✅ Checksum created: p3ipam-$VERSION-darwin-amd64.tar.gz.sha256"

echo ""
echo "🎯 Next Steps:"
echo "1. Create a Git tag: git tag -a $VERSION -m 'Release $VERSION'"
echo "2. Push the tag: git push origin $VERSION"
echo "3. Go to GitHub: https://github.com/palmar/p3ipam/releases/new"
echo "4. Select tag: $VERSION"
echo "5. Title: p3ipam $VERSION - Binary-Only Release"
echo "6. Description: Copy from RELEASE_NOTES.md"
echo "7. Upload: p3ipam-$VERSION-darwin-amd64.tar.gz"
echo "8. Publish release!"
echo ""

# Create release notes
cat > RELEASE_NOTES.md << EOF
# p3ipam $VERSION - Binary-Only Release

## 🔧 Release Fix

This release corrects the previous v0.1.0 release to provide a clean, binary-only distribution.

## ✨ Features

- **Core IPAM Functionality**: Manage subnets, hosts, and discoveries
- **SQLite Database**: Lightweight, file-based storage with automatic initialization
- **Flexible Parent References**: Reference subnets by name, ID, or CIDR notation
- **Smart Display**: Shows meaningful subnet names instead of cryptic IDs
- **Comprehensive Search**: Search across all objects with unified query interface
- **Table Formatting**: Clean, readable output for large datasets
- **Environment Configuration**: Configurable database location via P3IPAM_DATADIR
- **List Subnet Functionality**: Drill down into specific networks to see hosts

## 🚀 Quick Start

\`\`\`bash
# Initialize database
./p3ipam init

# Add your first subnet
./p3ipam add subnet --cidr 192.168.1.0/24 --name home-network

# Add hosts
./p3ipam add host --parent home-network --address 192.168.1.1 --name router
./p3ipam add host --parent home-network --address 192.168.1.10 --name main-pc

# List and search
./p3ipam list subnets
./p3ipam list subnet home-network
./p3ipam search 192.168.1
\`\`\`

## 📋 System Requirements

- **OS**: macOS, Linux, Windows
- **Architecture**: AMD64
- **Dependencies**: None (statically linked Go binary)
- **Storage**: Minimal (SQLite database file)

## 🔧 Installation

1. Download the release archive
2. Extract: \`tar -xzf p3ipam-$VERSION-darwin-amd64.tar.gz\`
3. Run: \`./p3ipam help\`

## 📁 Files Included

- \`p3ipam\` - Executable binary
- \`schema.sql\` - Database schema

## 🔄 Changes from v0.1.0

- **Fixed**: Release archive now contains only the binary and essential files
- **Improved**: Cleaner distribution without source code files
- **Optimized**: Smaller archive size for faster downloads

## 🐛 Known Issues

- None at this time

## 🔮 Future Plans

- IPv4 ping functionality with host discovery
- Delete and edit operations for objects
- Enhanced search filters
- Data export/import capabilities
- IPv6 support

## 📝 License

This project is released under the [Unlicense](LICENSE), placing it in the public domain.

**Note**: This repository contains code created using generative AI.

---

**p3ipam** - Simple IP Address Management for the rest of us.
EOF

echo "📝 Release notes created: RELEASE_NOTES.md"
echo ""
echo "🎯 Ready to create GitHub release!"
echo "Run: git tag -a $VERSION -m 'Release $VERSION'"
