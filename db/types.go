package db

import "time"

// Subnet represents a network subnet
type Subnet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CIDR      string    `json:"cidr"`
	ParentID  *string   `json:"parent_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// Host represents a network host
type Host struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	ParentID  string     `json:"parent_id"`
	Comment   string     `json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	LastSeen  *time.Time `json:"last_seen"`
}

// Discovery represents a discovered host from ping
type Discovery struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	SubnetID     string    `json:"subnet_id"`
	DiscoveredAt time.Time `json:"discovered_at"`
	LastSeen     time.Time `json:"last_seen"`
	Status       string    `json:"status"`
}

// SearchResults contains search results from all tables
type SearchResults struct {
	Subnets     []Subnet    `json:"subnets"`
	Hosts       []Host      `json:"hosts"`
	Discoveries []Discovery `json:"discoveries"`
}
