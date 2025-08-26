package utils

import (
	"fmt"
	"strings"

	"p3ipam/db"
)

// Table represents a formatted table for terminal output
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable creates a new table with the given headers
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		widths:  widths,
	}
}

// AddRow adds a row to the table and updates column widths
func (t *Table) AddRow(cells ...string) {
	if len(cells) != len(t.headers) {
		// Pad with empty strings if row is too short
		for len(cells) < len(t.headers) {
			cells = append(cells, "")
		}
		// Truncate if row is too long
		if len(cells) > len(t.headers) {
			cells = cells[:len(t.headers)]
		}
	}
	
	// Update column widths
	for i, cell := range cells {
		if len(cell) > t.widths[i] {
			t.widths[i] = len(cell)
		}
	}
	
	t.rows = append(t.rows, cells)
}

// String returns the formatted table as a string
func (t *Table) String() string {
	if len(t.rows) == 0 {
		return "No data to display.\n"
	}
	
	var result strings.Builder
	
	// Print header separator
	result.WriteString(t.printSeparator())
	
	// Print headers
	result.WriteString(t.printRow(t.headers))
	
	// Print header separator
	result.WriteString(t.printSeparator())
	
	// Print rows
	for _, row := range t.rows {
		result.WriteString(t.printRow(row))
	}
	
	// Print footer separator
	result.WriteString(t.printSeparator())
	
	// Print row count
	result.WriteString(fmt.Sprintf("\nTotal: %d rows\n", len(t.rows)))
	
	return result.String()
}

// printSeparator prints a separator line
func (t *Table) printSeparator() string {
	var result strings.Builder
	result.WriteString("+")
	for _, width := range t.widths {
		result.WriteString(strings.Repeat("-", width+2))
		result.WriteString("+")
	}
	result.WriteString("\n")
	return result.String()
}

// printRow prints a single row with proper padding
func (t *Table) printRow(cells []string) string {
	var result strings.Builder
	result.WriteString("|")
	for i, cell := range cells {
		result.WriteString(" ")
		result.WriteString(cell)
		result.WriteString(strings.Repeat(" ", t.widths[i]-len(cell)))
		result.WriteString(" |")
	}
	result.WriteString("\n")
	return result.String()
}

// FormatSubnets formats subnet data into a table
func FormatSubnets(subnets []db.Subnet) string {
	table := NewTable("ID", "CIDR", "Name", "Parent", "Comment", "Created")
	
	for _, subnet := range subnets {
		parent := ""
		if subnet.ParentID != nil {
			parent = *subnet.ParentID
		}
		
		table.AddRow(
			subnet.ID,
			subnet.CIDR,
			subnet.Name,
			parent,
			subnet.Comment,
			subnet.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
	
	return table.String()
}

// FormatHosts formats host data into a table
func FormatHosts(hosts []db.Host, subnetNames map[string]string) string {
	table := NewTable("ID", "Address", "Name", "Parent", "Comment", "Created", "Last Seen")
	
	for _, host := range hosts {
		parent := host.ParentID
		if name, exists := subnetNames[host.ParentID]; exists && name != "" {
			parent = name
		}
		
		lastSeen := ""
		if host.LastSeen != nil {
			lastSeen = host.LastSeen.Format("2006-01-02 15:04")
		}
		
		table.AddRow(
			host.ID,
			host.Address,
			host.Name,
			parent,
			host.Comment,
			host.CreatedAt.Format("2006-01-02 15:04"),
			lastSeen,
		)
	}
	
	return table.String()
}

// FormatDiscoveries formats discovery data into a table
func FormatDiscoveries(discoveries []db.Discovery, subnetNames map[string]string) string {
	table := NewTable("ID", "Address", "Subnet", "Status", "Discovered", "Last Seen")
	
	for _, discovery := range discoveries {
		subnet := discovery.SubnetID
		if name, exists := subnetNames[discovery.SubnetID]; exists && name != "" {
			subnet = name
		}
		
		table.AddRow(
			discovery.ID,
			discovery.Address,
			subnet,
			discovery.Status,
			discovery.DiscoveredAt.Format("2006-01-02 15:04"),
			discovery.LastSeen.Format("2006-01-02 15:04"),
		)
	}
	
	return table.String()
}
