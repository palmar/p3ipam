package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

// GetDatabasePath returns the database path from environment variable or default
func GetDatabasePath() string {
	if datadir := os.Getenv("P3IPAM_DATADIR"); datadir != "" {
		return filepath.Join(datadir, "p3ipam.db")
	}
	return "/opt/p3ipam/.data/p3ipam.db"
}

// Generate a pretty 6-character alphanumeric ID
func generateID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		if i < 3 {
			// First 3 characters are letters
			result[i] = charset[rand.Intn(26)]
		} else {
			// Last 3 characters are numbers
			result[i] = charset[26+rand.Intn(10)]
		}
	}
	return string(result)
}

// Connect to SQLite database
func Connect(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Connect to database
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := &Database{conn: conn}
	return db, nil
}

// Initialize database with schema
func (db *Database) Init() error {
	// Try to find schema.sql in multiple locations
	schemaPaths := []string{
		"schema.sql",         // Current directory
		"./schema.sql",       // Current directory explicit
		"../schema.sql",      // Parent directory
		"./db/../schema.sql", // Relative to db package
	}

	var schema []byte
	var err error

	for _, path := range schemaPaths {
		schema, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to read schema.sql from any location: %v", err)
	}

	_, err = db.conn.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	return nil
}

// Close database connection
func (db *Database) Close() error {
	return db.conn.Close()
}

// Get a unique ID that doesn't exist in any table
func (db *Database) GetUniqueID() string {
	for {
		id := generateID()

		// Check if ID exists in any table
		var count int

		// Check subnets
		err := db.conn.QueryRow("SELECT COUNT(*) FROM subnets WHERE id = ?", id).Scan(&count)
		if err != nil || count > 0 {
			continue
		}

		// Check hosts
		err = db.conn.QueryRow("SELECT COUNT(*) FROM hosts WHERE id = ?", id).Scan(&count)
		if err != nil || count > 0 {
			continue
		}

		// Check discoveries
		err = db.conn.QueryRow("SELECT COUNT(*) FROM discoveries WHERE id = ?", id).Scan(&count)
		if err != nil || count > 0 {
			continue
		}

		return id
	}
}

// Search across all tables
func (db *Database) Search(query string) (*SearchResults, error) {
	results := &SearchResults{}

	// Search subnets
	subnets, err := db.searchSubnets(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search subnets: %v", err)
	}
	results.Subnets = subnets

	// Search hosts
	hosts, err := db.searchHosts(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search hosts: %v", err)
	}
	results.Hosts = hosts

	// Search discoveries
	discoveries, err := db.searchDiscoveries(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search discoveries: %v", err)
	}
	results.Discoveries = discoveries

	return results, nil
}

func (db *Database) searchSubnets(query string) ([]Subnet, error) {
	rows, err := db.conn.Query(`
		SELECT id, name, cidr, parent_id, comment, created_at 
		FROM subnets 
		WHERE cidr LIKE ? OR name LIKE ? OR comment LIKE ?
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subnets []Subnet
	for rows.Next() {
		var s Subnet
		err := rows.Scan(&s.ID, &s.Name, &s.CIDR, &s.ParentID, &s.Comment, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, s)
	}

	return subnets, nil
}

func (db *Database) searchHosts(query string) ([]Host, error) {
	rows, err := db.conn.Query(`
		SELECT id, name, address, parent_id, comment, created_at, last_seen 
		FROM hosts 
		WHERE address LIKE ? OR name LIKE ? OR comment LIKE ?
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []Host
	for rows.Next() {
		var h Host
		err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.ParentID, &h.Comment, &h.CreatedAt, &h.LastSeen)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}

	return hosts, nil
}

func (db *Database) searchDiscoveries(query string) ([]Discovery, error) {
	rows, err := db.conn.Query(`
		SELECT id, address, subnet_id, discovered_at, last_seen, status 
		FROM discoveries 
		WHERE address LIKE ? OR status LIKE ?
	`, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discoveries []Discovery
	for rows.Next() {
		var d Discovery
		err := rows.Scan(&d.ID, &d.Address, &d.SubnetID, &d.DiscoveredAt, &d.LastSeen, &d.Status)
		if err != nil {
			return nil, err
		}
		discoveries = append(discoveries, d)
	}

	return discoveries, nil
}

// AddSubnet adds a new subnet to the database
func (db *Database) AddSubnet(cidr, name, parentRef, comment string) (*Subnet, error) {
	id := db.GetUniqueID()

	var parentIDPtr *string
	if parentRef != "" {
		// Resolve parent reference (name, ID, or CIDR)
		parentID, err := db.ResolveParentReference(parentRef)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parent reference '%s': %v", parentRef, err)
		}
		parentIDPtr = &parentID
	}

	_, err := db.conn.Exec(`
		INSERT INTO subnets (id, name, cidr, parent_id, comment, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, id, name, cidr, parentIDPtr, comment)

	if err != nil {
		return nil, fmt.Errorf("failed to insert subnet: %v", err)
	}

	subnet := &Subnet{
		ID:        id,
		Name:      name,
		CIDR:      cidr,
		ParentID:  parentIDPtr,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	return subnet, nil
}

// AddHost adds a new host to the database
func (db *Database) AddHost(address, name, parentRef, comment string) (*Host, error) {
	id := db.GetUniqueID()

	var parentID string
	if parentRef != "" {
		// Resolve parent reference (name, ID, or CIDR)
		var err error
		parentID, err = db.ResolveParentReference(parentRef)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parent reference '%s': %v", parentRef, err)
		}
	}

	_, err := db.conn.Exec(`
		INSERT INTO hosts (id, name, address, parent_id, comment, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, id, name, address, parentID, comment)

	if err != nil {
		return nil, fmt.Errorf("failed to insert host: %v", err)
	}

	host := &Host{
		ID:        id,
		Name:      name,
		Address:   address,
		ParentID:  parentID,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	return host, nil
}

// ResolveParentReference resolves a parent reference by name, ID, or CIDR
// Returns the parent ID and an error if there are multiple matches
func (db *Database) ResolveParentReference(reference string) (string, error) {
	if reference == "" {
		return "", nil
	}

	// Search for matches in all three fields
	var matches []string

	// Check by ID (exact match)
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM subnets WHERE id = ?", reference).Scan(&count)
	if err == nil && count > 0 {
		matches = append(matches, reference) // ID is unique, so this is the match
	}

	// Check by name (exact match)
	err = db.conn.QueryRow("SELECT COUNT(*) FROM subnets WHERE name = ?", reference).Scan(&count)
	if err == nil && count > 0 {
		matches = append(matches, reference)
	}

	// Check by CIDR (exact match)
	err = db.conn.QueryRow("SELECT COUNT(*) FROM subnets WHERE cidr = ?", reference).Scan(&count)
	if err == nil && count > 0 {
		matches = append(matches, reference)
	}

	// Handle results
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no subnet found matching reference: %s", reference)
	case 1:
		// Get the actual ID for the match
		var id string
		if err := db.conn.QueryRow("SELECT id FROM subnets WHERE id = ? OR name = ? OR cidr = ?", reference, reference, reference).Scan(&id); err != nil {
			return "", fmt.Errorf("failed to get parent ID: %v", err)
		}
		return id, nil
	default:
		return "", fmt.Errorf("multiple subnets match reference '%s'. Please use a more specific reference (ID, unique name, or exact CIDR)", reference)
	}
}

// Initialize random seed for ID generation
func init() {
	rand.Seed(time.Now().UnixNano())
}
