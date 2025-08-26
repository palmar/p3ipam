# p3ipam - Lightweight IP Address Management Tool

A simple, lightweight terminal application for managing IP addresses, subnets, and network hosts. Built in Go with SQLite backend.

## Features

- **Simple Command Structure**: `p3ipam <action> <object> <target> [--arguments]`
- **Pretty 6-Character IDs**: Easy-to-read identifiers (e.g., ABC123)
- **Unified Search**: Search across all objects with a single query
- **SQLite Database**: Lightweight, file-based storage
- **IPv4 Ping Support**: Discover live hosts on networks

## Installation

1. Ensure you have Go 1.21+ installed
2. Clone this repository
3. Run `go mod tidy` to download dependencies
4. Build with `go build -o p3ipam main.go`

## Quick Start

1. **Initialize the database**:
   ```bash
   p3ipam init
   ```

2. **Add a subnet**:
   ```bash
   p3ipam add subnet --cidr 192.168.1.0/24 --name home-network
   ```

3. **Add a host**:
   ```bash
   p3ipam add host --parent home-network --address 192.168.1.1 --name router
   ```

4. **Search for anything**:
   ```bash
   p3ipam search 192.168.1
   ```

5. **Ping a subnet to discover hosts**:
   ```bash
   p3ipam ping subnet home-network
   ```

## Commands

### Actions
- `init` - Initialize the database
- `version` - Show version information
- `help` - Show help message
- `add` - Add a new object
- `list` - List objects
- `delete` - Delete an object
- `edit` - Edit an object
- `ping` - Ping and discover hosts
- `search` - Search across all objects

### Objects
- `subnet` - Network subnet (e.g., 192.168.1.0/24)
- `host` - Network host (e.g., 192.168.1.1)

## Database Schema

The tool uses SQLite with three main tables:

- **subnets**: Network subnets with CIDR notation
- **hosts**: Individual network hosts
- **discoveries**: Results from ping sweeps

## Current Status

This is a **Minimum Viable Product (MVP)** with:
- ✅ Database initialization and schema
- ✅ Basic command structure
- ✅ Search functionality across all tables
- ✅ Pretty 6-character ID system
- ✅ Help and version commands
- ✅ Error handling

**Coming Soon:**
- Add/Edit/Delete operations for subnets and hosts
- IPv4 ping functionality with host discovery
- Enhanced search with filters
- Data export/import capabilities

## Development

The project structure:
```
p3ipam/
├── main.go          # Main application and CLI handling
├── db/
│   ├── database.go  # Database operations
│   └── types.go     # Data structures
├── schema.sql       # Database schema
└── go.mod           # Go module definition
```

## License

This is a personal project for learning and network management.
