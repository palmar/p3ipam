-- p3ipam Database Schema

-- Subnets table
CREATE TABLE IF NOT EXISTS subnets (
    id TEXT PRIMARY KEY,           -- 6-char alphanumeric ID (e.g., ABC123)
    name TEXT,                     -- Optional name
    cidr TEXT NOT NULL,            -- CIDR notation (e.g., 192.168.1.0/24)
    parent_id TEXT,                -- Parent subnet ID (NULL for root)
    comment TEXT,                  -- Optional comment
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES subnets(id)
);

-- Hosts table
CREATE TABLE IF NOT EXISTS hosts (
    id TEXT PRIMARY KEY,           -- 6-char alphanumeric ID (e.g., ABC123)
    name TEXT,                     -- Optional name
    address TEXT NOT NULL,         -- IP address
    parent_id TEXT,                -- Parent subnet ID
    comment TEXT,                  -- Optional comment
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME,            -- When host was last pinged
    FOREIGN KEY (parent_id) REFERENCES subnets(id)
);

-- Discoveries table (for ping results)
CREATE TABLE IF NOT EXISTS discoveries (
    id TEXT PRIMARY KEY,           -- 6-char alphanumeric ID (e.g., ABC123)
    address TEXT NOT NULL,         -- IP address discovered
    subnet_id TEXT,                -- Subnet where it was discovered
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'alive',   -- alive, dead, unknown
    FOREIGN KEY (subnet_id) REFERENCES subnets(id)
);

-- Indexes for better search performance
CREATE INDEX IF NOT EXISTS idx_subnets_cidr ON subnets(cidr);
CREATE INDEX IF NOT EXISTS idx_subnets_name ON subnets(name);
CREATE INDEX IF NOT EXISTS idx_hosts_address ON hosts(address);
CREATE INDEX IF NOT EXISTS idx_hosts_name ON hosts(name);
CREATE INDEX IF NOT EXISTS idx_discoveries_address ON discoveries(address);
CREATE INDEX IF NOT EXISTS idx_discoveries_subnet ON discoveries(subnet_id);
