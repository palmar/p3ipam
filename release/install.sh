#!/bin/bash

# p3ipam Installation Script
# This script installs p3ipam to /usr/local/bin

set -e

echo "Installing p3ipam..."

# Check if binary exists
if [ ! -f "./p3ipam" ]; then
    echo "Error: p3ipam binary not found in current directory"
    echo "Please run this script from the release directory"
    exit 1
fi

# Install binary
echo "Installing binary to /usr/local/bin..."
sudo cp p3ipam /usr/local/bin/
sudo chmod +x /usr/local/bin/p3ipam

# Create data directory
echo "Creating default data directory..."
sudo mkdir -p /opt/p3ipam/.data
sudo chown $USER:$USER /opt/p3ipam/.data

echo ""
echo "✅ p3ipam installed successfully!"
echo ""
echo "Next steps:"
echo "1. Initialize the database: p3ipam init"
echo "2. Add your first subnet: p3ipam add subnet --cidr 192.168.1.0/24 --name home-network"
echo "3. Run 'p3ipam help' for more information"
echo ""
echo "Note: The default database location is /opt/p3ipam/.data/"
echo "You can override this with: export P3IPAM_DATADIR=/path/to/your/database"
