# p3ipam - Lightweight IP Address Management Tool

A simple, fast terminal-based IP Address Management (IPAM) tool built in Go. Perfect for homelabs, small offices, and network administrators who need a lightweight solution without the complexity of enterprise tools.

## Features

- **Simple CLI Interface**: Easy-to-use command structure
- **SQLite Database**: Lightweight, file-based storage
- **Flexible Parent References**: Reference subnets by name, ID, or CIDR
- **Smart Display**: Shows meaningful names instead of cryptic IDs
- **Comprehensive Search**: Search across all objects with one command
- **Table Formatting**: Clean, readable output for large datasets
- **Environment Configuration**: Configurable database location via environment variable

## Quick Start

### 1. Initialize the Database
```bash
./p3ipam init
```
This will prompt for:
- Database location (default: `/opt/p3ipam/.data/`)
- Database password (default: empty)

**Note**: If you choose a custom database location, set the environment variable:
```bash
export P3IPAM_DATADIR=/path/to/your/database
```

### 2. Add Your First Subnet
```bash
./p3ipam add subnet --cidr 192.168.1.0/24 --name home-network --comment "Home network for family devices"
```

### 3. Add Hosts
```bash
./p3ipam add host --parent home-network --address 192.168.1.1 --name router --comment "Old Asus router lol"
./p3ipam add host --parent home-network --address 192.168.1.10 --name main-pc --comment "Workstation for me"
```

### 4. List and Search
```bash
./p3ipam list subnets
./p3ipam list hosts
./p3ipam list subnet home-network
./p3ipam search 192.168.1
```

## Command Reference

### Basic Commands
- `p3ipam init` - Initialize the database
- `p3ipam version` - Show version information
- `p3ipam help` - Show help message

### Adding Objects
- `p3ipam add subnet --cidr <CIDR> [--name <name>] [--parent <reference>] [--comment <comment>]`
- `p3ipam add host --address <IP> [--name <name>] --parent <reference> [--comment <comment>]`

### Listing Objects
- `p3ipam list subnets` - List all subnets
- `p3ipam list hosts` - List all hosts
- `p3ipam list discoveries` - List all discoveries
- `p3ipam list subnet <reference>` - List hosts in a specific subnet

### Searching
- `p3ipam search <query>` - Search across all objects

### Parent References
The `--parent` argument accepts:
- **Subnet name**: `--parent home-network`
- **Subnet ID**: `--parent ABC123`
- **CIDR notation**: `--parent 192.168.1.0/24`

## Examples

### Network Documentation
```bash
# Create a home network
./p3ipam add subnet --cidr 192.168.1.0/24 --name home-network --comment "Home network for family devices"

# Add network devices
./p3ipam add host --parent home-network --address 192.168.1.1 --name router --comment "Asus RT-AC68U"
./p3ipam add host --parent home-network --address 192.168.1.10 --name main-pc --comment "Gaming PC"
./p3ipam add host --parent home-network --address 192.168.1.15 --name nas --comment "Synology DS920+"
./p3ipam add host --parent home-network --address 192.168.1.20 --name laptop --comment "MacBook Pro"

# Create a guest network
./p3ipam add subnet --cidr 192.168.2.0/24 --name guest-network --comment "Guest WiFi network"
./p3ipam add host --parent guest-network --address 192.168.2.1 --name guest-router --comment "Guest network gateway"
```

### Office Network
```bash
# Main office network
./p3ipam add subnet --cidr 10.0.1.0/24 --name office-main --comment "Main office network - VLAN 10"

# Add office devices
./p3ipam add host --parent office-main --address 10.0.1.1 --name gateway --comment "Cisco ASA 5506-X"
./p3ipam add host --parent office-main --address 10.0.1.10 --name dhcp-server --comment "Windows Server 2019 DHCP"
./p3ipam add host --parent office-main --address 10.0.1.100 --name file-server --comment "Dell PowerEdge R740"

# Server network
./p3ipam add subnet --cidr 10.0.2.0/24 --name office-servers --comment "Server network - VLAN 20"
./p3ipam add host --parent office-servers --address 10.0.2.10 --name web-server --comment "Nginx web server"
./p3ipam add host --parent office-servers --address 10.0.2.20 --name db-server --comment "PostgreSQL database"
```

## Database Schema

The tool uses SQLite with the following structure:
- **subnets**: Network subnets with CIDR notation
- **hosts**: Individual hosts with IP addresses
- **discoveries**: Discovered hosts from ping operations

## Environment Variables

- `P3IPAM_DATADIR`: Override default database location

## Requirements

- **OS**: Linux, macOS, Windows
- **Dependencies**: None (statically linked Go binary)
- **Storage**: Minimal (SQLite database file)

## Building from Source

```bash
git clone <repository>
cd p3ipam
go build -o p3ipam main.go
```

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

---

**p3ipam** - Simple IP Address Management for the rest of us.
