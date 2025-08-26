# p3ipam

A lightweight, terminal-based IP Address Management (IPAM) tool built in Go. Perfect for homelabs, small offices, and network administrators who need a simple solution without enterprise complexity.

## Features

- **Simple CLI Interface**: `p3ipam <action> <object> <target> [--arguments]`
- **SQLite Database**: Lightweight, file-based storage
- **Flexible Parent References**: Reference subnets by name, ID, or CIDR
- **Smart Display**: Shows meaningful names instead of cryptic IDs
- **Comprehensive Search**: Search across all objects with one command
- **Table Formatting**: Clean, readable output for large datasets
- **Environment Configuration**: Configurable database location

## Quick Start

```bash
# Initialize database
p3ipam init

# Add a subnet
p3ipam add subnet --cidr 192.168.1.0/24 --name home-network

# Add hosts
p3ipam add host --parent home-network --address 192.168.1.1 --name router
p3ipam add host --parent home-network --address 192.168.1.10 --name main-pc

# List and search
p3ipam list subnets
p3ipam list subnet home-network
p3ipam search 192.168.1
```

## Installation

### From Release
Download the latest release from the [Releases](https://github.com/palmarg/p3ipam/releases) page.

### From Source
```bash
git clone https://github.com/palmarg/p3ipam.git
cd p3ipam
go build -o p3ipam main.go
```

## Documentation

See the [release README](release/README.md) for comprehensive documentation and examples.

## License

This project is released under the [Unlicense](LICENSE), placing it in the public domain.

## Contributing

[Add contribution guidelines here]

## AI Disclosure

**Note**: This repository contains code created using generative AI.
